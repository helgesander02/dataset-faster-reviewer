package cache

import (
	"backend/src/models"
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"github.com/patrickmn/go-cache"
)

type Base64Cache struct {
	cache *cache.Cache
}

func NewBase64Cache() *Base64Cache {
	// 設定 cache 過期時間為 30 分鐘，清理間隔為 10 分鐘
	c := cache.New(30*time.Minute, 10*time.Minute)
	return &Base64Cache{
		cache: c,
	}
}

func (bc *Base64Cache) GetAllPageDetail(jobName string) map[int]string {
	log.Printf("Fetching all page details from cache for job: %s", jobName)
	
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		log.Printf("No cache found for job: %s", jobName)
		return nil
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		log.Printf("Invalid cache data type for job: %s", jobName)
		return nil
	}

	detail := make(map[int]string)
	for idx, page := range imagesPerPage.Pages {
		detail[idx] = page.DatasetName
	}
	return detail
}

func (bc *Base64Cache) InitialImagesCache(jobName string, pageNumber int, getDatasetsFn func(string) []string, getImagesFn func(string, string) []models.Image) {
	log.Printf("Initializing images cache for job: %s with page size: %d", jobName, pageNumber)
	
	imagesPerPage := models.NewImagesPerPageCache()
	imagesPerPage.JobName = jobName

	allDatasets := getDatasetsFn(jobName)
	for _, datasetName := range allDatasets {
		images := getImagesFn(jobName, datasetName)
		for i := 0; i < len(images); i += pageNumber {
			end := i + pageNumber
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
}

func (bc *Base64Cache) GetImagesCache(jobName string, pageIndex int, pageNumber int, getDatasetsFn func(string) []string, getImagesFn func(string, string) []models.Image) (models.ImageItems, bool) {
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		log.Printf("Cache not found for job: %s, initializing...", jobName)
		bc.InitialImagesCache(jobName, pageNumber, getDatasetsFn, getImagesFn)
		cacheData, found = bc.cache.Get(jobName)
		if !found {
			log.Printf("Failed to initialize cache for job: %s", jobName)
			return models.ImageItems{}, false
		}
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		log.Printf("Invalid cache data type for job: %s", jobName)
		return models.ImageItems{}, false
	}

	if len(imagesPerPage.Pages) <= pageIndex {
		log.Printf("Page index %d out of range for job %s", pageIndex, jobName)
		return models.ImageItems{}, false
	}

	cachedImages := imagesPerPage.Pages[pageIndex]
	
	// 如果 Base64ImageSet 為空，則生成 base64 圖片
	if len(cachedImages.Base64ImageSet) == 0 {
		log.Printf("Generating base64 images for job %s at page index %d", jobName, pageIndex)
		bc.generateBase64Images(jobName, pageIndex)
		
		// 重新獲取更新後的數據
		cacheData, _ = bc.cache.Get(jobName)
		imagesPerPage = cacheData.(models.ImagesPerPageCache)
		cachedImages = imagesPerPage.Pages[pageIndex]
	}

	exist := len(cachedImages.Base64ImageSet) > 0
	return cachedImages, exist
}

func (bc *Base64Cache) generateBase64Images(jobName string, pageIndex int) {
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		log.Printf("Cache not found for job: %s", jobName)
		return
	}

	imagesPerPage, ok := cacheData.(models.ImagesPerPageCache)
	if !ok {
		log.Printf("Invalid cache data type for job: %s", jobName)
		return
	}

	if pageIndex >= len(imagesPerPage.Pages) {
		log.Printf("Page index %d out of range for job %s", pageIndex, jobName)
		return
	}

	images := imagesPerPage.Pages[pageIndex].ImageSet
	for _, img := range images {
		base64Img := compressImages(img.Path)
		if base64Img != "" {
			compressedImage := models.NewImage(img.Name, base64Img)
			imagesPerPage.Pages[pageIndex].Base64ImageSet = append(imagesPerPage.Pages[pageIndex].Base64ImageSet, compressedImage)
		}
	}

	// 更新 cache
	bc.cache.Set(jobName, imagesPerPage, cache.DefaultExpiration)
	log.Printf("Generated %d base64 images for job %s at page index %d", len(imagesPerPage.Pages[pageIndex].Base64ImageSet), jobName, pageIndex)
}

func (bc *Base64Cache) ClearImagesCache(jobName string) {
	bc.cache.Delete(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
}

func (bc *Base64Cache) ImageCacheJobExists(jobName string) bool {
	_, found := bc.cache.Get(jobName)
	return found
}

func compressImages(imgPath string) string {
	log.Println("Processing image:", imgPath)
	file, err := os.Open(imgPath)
	if err != nil {
		log.Println("Failed to open image:", err)
		return ""
	}
	defer file.Close()

	decodedImg, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Failed to decode image (format: %s): %v\n", format, err)
		return ""
	}

	resizedImg := resize.Resize(150, 0, decodedImg, resize.Lanczos3)

	var buf bytes.Buffer
	opts := &webp.Options{Lossless: false, Quality: 75}
	if err := webp.Encode(&buf, resizedImg, opts); err != nil {
		log.Println("Failed to encode WebP:", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
