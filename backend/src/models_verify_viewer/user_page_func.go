// user_page_func.go
package models_verify_viewer

import (
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type PageCache struct {
	cache          *cache.Cache
	dataProvider   DataProvider
	imageProcessor ImageProcessor
}

func NewPageCache(dataProvider DataProvider, imageProcessor ImageProcessor) *PageCache {
	return &PageCache{
		cache:          cache.New(30*time.Minute, 10*time.Minute),
		dataProvider:   dataProvider,
		imageProcessor: imageProcessor,
	}
}

func (pc *PageCache) GetAllPageDetail(jobName string) (map[int]string, error) {
	log.Printf("Fetching all page details from cache for job: %s", jobName)

	cacheData, found := pc.cache.Get(jobName)
	if !found {
		return nil, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
	if !ok {
		pc.cache.Delete(jobName) // 清理無效數據
		return nil, ErrInvalidCacheData{JobName: jobName}
	}

	detail := make(map[int]string)
	for idx, page := range imagesPerPage.Pages {
		detail[idx] = page.DatasetName
	}
	return detail, nil
}

func (pc *PageCache) InitializeImagesCache(jobName string, pageSize int) error {
	log.Printf("Initializing images cache for job: %s with page size: %d", jobName, pageSize)

	imagesPerPage := NewImagesPerPageCache()
	imagesPerPage.JobName = jobName

	allDatasets := pc.dataProvider.GetDatasets(jobName)
	for _, datasetName := range allDatasets {
		images := pc.dataProvider.GetImages(jobName, datasetName)

		// 分頁處理
		for i := 0; i < len(images); i += pageSize {
			end := i + pageSize
			if end > len(images) {
				end = len(images)
			}

			imageSet := ImageItems{
				DatasetName:    datasetName,
				ImageSet:       images[i:end],
				Base64ImageSet: []Image{},
			}
			imagesPerPage.Pages = append(imagesPerPage.Pages, imageSet)
		}
	}

	imagesPerPage.MaxPage = len(imagesPerPage.Pages)
	pc.cache.Set(jobName, imagesPerPage, cache.DefaultExpiration)

	log.Printf("Initialized images cache for job %s with %d pages", jobName, imagesPerPage.MaxPage)
	return nil
}

func (pc *PageCache) GetImagesCache(jobName string, pageIndex int) (ImageItems, error) {
	cacheData, found := pc.cache.Get(jobName)
	if !found {
		return ImageItems{}, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
	if !ok {
		pc.cache.Delete(jobName)
		return ImageItems{}, ErrInvalidCacheData{JobName: jobName}
	}

	if pageIndex >= len(imagesPerPage.Pages) || pageIndex < 0 {
		return ImageItems{}, ErrPageIndexOutOfRange{
			JobName:   jobName,
			PageIndex: pageIndex,
			MaxPages:  len(imagesPerPage.Pages),
		}
	}

	cachedImages := imagesPerPage.Pages[pageIndex]

	// 懶加載 base64 圖片
	if len(cachedImages.Base64ImageSet) == 0 {
		if err := pc.generateBase64Images(jobName, pageIndex); err != nil {
			log.Printf("Failed to generate base64 images: %v", err)
			return cachedImages, err
		}

		// 重新獲取更新後的數據
		cacheData, _ = pc.cache.Get(jobName)
		imagesPerPage = cacheData.(ImagesPerPageCache)
		cachedImages = imagesPerPage.Pages[pageIndex]
	}

	return cachedImages, nil
}

func (pc *PageCache) generateBase64Images(jobName string, pageIndex int) error {
	cacheData, found := pc.cache.Get(jobName)
	if !found {
		return ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
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
	var processedImages []Image

	for _, img := range images {
		base64Img, err := pc.imageProcessor.CompressToBase64(img.Path)
		if err != nil {
			log.Printf("Failed to process image %s: %v", img.Path, err)
			continue // 跳過有問題的圖片，繼續處理其他圖片
		}

		if base64Img != "" {
			processedImages = append(processedImages, NewImage(img.Name, base64Img))
		}
	}

	imagesPerPage.Pages[pageIndex].Base64ImageSet = processedImages
	pc.cache.Set(jobName, imagesPerPage, cache.DefaultExpiration)

	log.Printf("Generated %d base64 images for job %s at page index %d",
		len(processedImages), jobName, pageIndex)
	return nil
}

func (pc *PageCache) ClearImagesCache(jobName string) {
	pc.cache.Delete(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
}

func (pc *PageCache) ImageCacheJobExists(jobName string) bool {
	_, found := pc.cache.Get(jobName)
	return found
}

func (pc *PageCache) GetMaxPages(jobName string) (int, error) {
	cacheData, found := pc.cache.Get(jobName)
	if !found {
		return 0, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
	if !ok {
		return 0, ErrInvalidCacheData{JobName: jobName}
	}

	return imagesPerPage.MaxPage, nil
}
