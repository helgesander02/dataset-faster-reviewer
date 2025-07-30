package cache

import (
	"backend/src/models"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type Base64Cache struct {
	cache          *cache.Cache
	dataProvider   DataProvider
	imageProcessor ImageProcessor
}

func NewBase64Cache(dataProvider DataProvider, imageProcessor ImageProcessor) *Base64Cache {
	c := cache.New(30*time.Minute, 10*time.Minute)
	return &Base64Cache{
		cache:          c,
		dataProvider:   dataProvider,
		imageProcessor: imageProcessor,
	}
}

func (bc *Base64Cache) GetAllPageDetail(jobName string) (map[int]string, error) {
	log.Printf("Fetching all page details from cache for job: %s", jobName)
	
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		return nil, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		bc.cache.Delete(jobName) // 清理無效數據
		return nil, ErrInvalidCacheData{JobName: jobName}
	}

	detail := make(map[int]string)
	for idx, page := range imagesPerPage.Pages {
		detail[idx] = page.DatasetName
	}
	return detail, nil
}

func (bc *Base64Cache) InitializeImagesCache(jobName string, pageSize int) error {
	log.Printf("Initializing images cache for job: %s with page size: %d", jobName, pageSize)
	
	imagesPerPage := models.NewImagesPerPageCache()
	imagesPerPage.JobName = jobName

	allDatasets := bc.dataProvider.GetDatasets(jobName)
	for _, datasetName := range allDatasets {
		images := bc.dataProvider.GetImages(jobName, datasetName)
		
		// 分頁處理
		for i := 0; i < len(images); i += pageSize {
			end := i + pageSize
			if end > len(images) {
				end = len(images)
			}
			
			imageSet := models.ImageItems{
				DatasetName:    datasetName,
				ImageSet:       images[i:end],
				Base64ImageSet: []models.Image{},
			}
			imagesPerPage.Pages = append(imagesPerPage.Pages, imageSet)
		}
	}
	
	imagesPerPage.MaxPage = len(imagesPerPage.Pages)
	bc.cache.Set(jobName, imagesPerPage, cache.DefaultExpiration)
	
	log.Printf("Initialized images cache for job %s with %d pages", jobName, imagesPerPage.MaxPage)
	return nil
}

func (bc *Base64Cache) GetImagesCache(jobName string, pageIndex int) (models.ImageItems, error) {
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		return models.ImageItems{}, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		bc.cache.Delete(jobName)
		return models.ImageItems{}, ErrInvalidCacheData{JobName: jobName}
	}

	if pageIndex >= len(imagesPerPage.Pages) || pageIndex < 0 {
		return models.ImageItems{}, ErrPageIndexOutOfRange{
			JobName:   jobName,
			PageIndex: pageIndex,
			MaxPages:  len(imagesPerPage.Pages),
		}
	}

	cachedImages := imagesPerPage.Pages[pageIndex]
	
	// 懶加載 base64 圖片
	if len(cachedImages.Base64ImageSet) == 0 {
		if err := bc.generateBase64Images(jobName, pageIndex); err != nil {
			log.Printf("Failed to generate base64 images: %v", err)
			return cachedImages, err
		}
		
		// 重新獲取更新後的數據
		cacheData, _ = bc.cache.Get(jobName)
		imagesPerPage = cacheData.(models.ImagesPerPageCache)
		cachedImages = imagesPerPage.Pages[pageIndex]
	}

	return cachedImages, nil
}

func (bc *Base64Cache) generateBase64Images(jobName string, pageIndex int) error {
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		return ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		return ErrInvalidCacheData{JobName: jobName}
	}

	if pageIndex >= len(imagesPerPage.Pages) {
		return ErrPageIndexOutOfRange{
			JobName:   jobName,
			PageIndex: pageIndex,
			MaxPages:  len(imagesPerPage.Pages),
		}
	}

	images := imagesPerPage.Pages[pageIndex].ImageSet
	var processedImages []models.Image
	
	for _, img := range images {
		base64Img, err := bc.imageProcessor.CompressToBase64(img.Path)
		if err != nil {
			log.Printf("Failed to process image %s: %v", img.Path, err)
			continue // 跳過有問題的圖片，繼續處理其他圖片
		}
		
		if base64Img != "" {
			processedImages = append(processedImages, models.NewImage(img.Name, base64Img))
		}
	}
	
	imagesPerPage.Pages[pageIndex].Base64ImageSet = processedImages
	bc.cache.Set(jobName, imagesPerPage, cache.DefaultExpiration)
	
	log.Printf("Generated %d base64 images for job %s at page index %d", 
		len(processedImages), jobName, pageIndex)
	return nil
}

func (bc *Base64Cache) ClearImagesCache(jobName string) {
	bc.cache.Delete(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
}

func (bc *Base64Cache) ImageCacheJobExists(jobName string) bool {
	_, found := bc.cache.Get(jobName)
	return found
}

func (bc *Base64Cache) GetMaxPages(jobName string) (int, error) {
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		return 0, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		return 0, ErrInvalidCacheData{JobName: jobName}
	}

	return imagesPerPage.MaxPage, nil
}
