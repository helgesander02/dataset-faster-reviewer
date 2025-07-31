package models_verify_viewer

import (
	"fmt"
	"log"

	"github.com/patrickmn/go-cache"
)

// Job Cache 相關方法
func (cm *CacheManager) GetJobCache(jobName string) (Job, bool) {
	return cm.JobCache.Get(jobName)
}

func (cm *CacheManager) SetJobCache(job Job) {
	cm.JobCache.Set(job)
}

func (cm *CacheManager) DeleteJobCache(jobName string) {
	cm.JobCache.Delete(jobName)
}

func (cm *CacheManager) JobCacheExists(jobName string) bool {
	return cm.JobCache.Exists(jobName)
}

func (cm *CacheManager) ClearJobCache() {
	cm.JobCache.Clear()
}
func (jc *JobCache) Get(jobName string) (Job, bool) {
	data, found := jc.cache.Get(jobName)
	if !found {
		return Job{}, false
	}

	job, ok := data.(Job)
	if !ok {
		log.Printf("Invalid job data type in cache for job: %s", jobName)
		jc.cache.Delete(jobName)
		return Job{}, false
	}

	log.Printf("Job found in cache: %s", jobName)
	return job, true
}

func (jc *JobCache) Set(job Job) {
	// 檢查是否超過最大項目數
	if jc.cache.ItemCount() >= jc.maxItems {
		// 刪除最舊的項目
		jc.evictOldest()
	}

	jc.cache.Set(job.Name, job, cache.DefaultExpiration)
	log.Printf("Added job to cache: %s", job.Name)
}

func (jc *JobCache) Delete(jobName string) {
	jc.cache.Delete(jobName)
	log.Printf("Removed job from cache: %s", jobName)
}

func (jc *JobCache) Exists(jobName string) bool {
	_, found := jc.cache.Get(jobName)
	return found
}

func (jc *JobCache) Clear() {
	jc.cache.Flush()
	log.Println("Cleared all job cache")
}

func (jc *JobCache) ItemCount() int {
	return jc.cache.ItemCount()
}

func (jc *JobCache) evictOldest() {
	items := jc.cache.Items()
	if len(items) == 0 {
		return
	}

	var oldestKey string
	var oldestExpiration int64
	first := true

	for key, item := range items {
		if first || item.Expiration < oldestExpiration {
			oldestKey = key
			oldestExpiration = item.Expiration
			first = false
		}
	}

	if oldestKey != "" {
		jc.cache.Delete(oldestKey)
		log.Printf("Evicted oldest job from cache: %s", oldestKey)
	}
}

// Base64 Cache 相關方法
func (cm *CacheManager) GetAllPageDetail(jobName string) (map[int]string, error) {
	return cm.Base64Cache.GetAllPageDetail(jobName)
}

func (cm *CacheManager) InitializeImagesCache(jobName string, pageSize int) error {
	return cm.Base64Cache.InitializeImagesCache(jobName, pageSize)
}

func (cm *CacheManager) GetImagesCache(jobName string, pageIndex int) (ImageItems, error) {
	return cm.Base64Cache.GetImagesCache(jobName, pageIndex)
}

func (cm *CacheManager) ClearImagesCache(jobName string) {
	cm.Base64Cache.ClearImagesCache(jobName)
}

func (cm *CacheManager) ImageCacheJobExists(jobName string) bool {
	return cm.Base64Cache.ImageCacheJobExists(jobName)
}

func (cm *CacheManager) GetJobMaxPages(jobName string) (int, error) {
	return cm.Base64Cache.GetMaxPages(jobName)
}

func (bc *Base64Cache) GetAllPageDetail(jobName string) (map[int]string, error) {
	log.Printf("Fetching all page details from cache for job: %s", jobName)

	cacheData, found := bc.cache.Get(jobName)
	if !found {
		return nil, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
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

	imagesPerPage := NewImagesPerPageCache()
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

			imageSet := ImageItems{
				DatasetName:    datasetName,
				ImageSet:       images[i:end],
				Base64ImageSet: []Image{},
			}
			imagesPerPage.Pages = append(imagesPerPage.Pages, imageSet)
		}
	}

	imagesPerPage.MaxPage = len(imagesPerPage.Pages)
	bc.cache.Set(jobName, imagesPerPage, cache.DefaultExpiration)

	log.Printf("Initialized images cache for job %s with %d pages", jobName, imagesPerPage.MaxPage)
	return nil
}

func (bc *Base64Cache) GetImagesCache(jobName string, pageIndex int) (ImageItems, error) {
	cacheData, found := bc.cache.Get(jobName)
	if !found {
		return ImageItems{}, ErrCacheNotFound{JobName: jobName}
	}

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
	if !ok {
		bc.cache.Delete(jobName)
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
		if err := bc.generateBase64Images(jobName, pageIndex); err != nil {
			log.Printf("Failed to generate base64 images: %v", err)
			return cachedImages, err
		}

		// 重新獲取更新後的數據
		cacheData, _ = bc.cache.Get(jobName)
		imagesPerPage = cacheData.(ImagesPerPageCache)
		cachedImages = imagesPerPage.Pages[pageIndex]
	}

	return cachedImages, nil
}

func (bc *Base64Cache) generateBase64Images(jobName string, pageIndex int) error {
	cacheData, found := bc.cache.Get(jobName)
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
		base64Img, err := bc.imageProcessor.CompressToBase64(img.Path)
		if err != nil {
			log.Printf("Failed to process image %s: %v", img.Path, err)
			continue // 跳過有問題的圖片，繼續處理其他圖片
		}

		if base64Img != "" {
			processedImages = append(processedImages, NewImage(img.Name, base64Img))
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

	imagesPerPage, ok := cacheData.(ImagesPerPageCache)
	if !ok {
		return 0, ErrInvalidCacheData{JobName: jobName}
	}

	return imagesPerPage.MaxPage, nil
}

// 清理所有快取
func (cm *CacheManager) ClearAllCache() {
	cm.JobCache.Clear()
	log.Println("Cleared all caches")
}

// 取得快取統計資訊
func (cm *CacheManager) GetCacheStats() map[string]int {
	return map[string]int{
		"job_cache_items": cm.JobCache.ItemCount(),
	}
}

type ErrCacheNotFound struct {
	JobName string
}

func (e ErrCacheNotFound) Error() string {
	return fmt.Sprintf("cache not found for job: %s", e.JobName)
}

type ErrInvalidCacheData struct {
	JobName string
}

func (e ErrInvalidCacheData) Error() string {
	return fmt.Sprintf("invalid cache data type for job: %s", e.JobName)
}

type ErrPageIndexOutOfRange struct {
	JobName   string
	PageIndex int
	MaxPages  int
}

func (e ErrPageIndexOutOfRange) Error() string {
	return fmt.Sprintf("page index %d out of range for job %s (max pages: %d)",
		e.PageIndex, e.JobName, e.MaxPages)
}
