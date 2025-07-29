package cache

import (
	"backend/src/models"
	"log"
)

type CacheManager struct {
	JobCache    *JobCache
	Base64Cache *Base64Cache
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		JobCache:    NewJobCache(3), // 最多快取 3 個 job
		Base64Cache: NewBase64Cache(),
	}
}

// Job Cache 相關方法
func (cm *CacheManager) GetJobCache(jobName string) (models.Job, bool) {
	return cm.JobCache.Get(jobName)
}

func (cm *CacheManager) SetJobCache(job models.Job) {
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

// Base64 Cache 相關方法
func (cm *CacheManager) GetAllPageDetail(jobName string) map[int]string {
	return cm.Base64Cache.GetAllPageDetail(jobName)
}

func (cm *CacheManager) InitialImagesCache(jobName string, pageNumber int, getDatasetsFn func(string) []string, getImagesFn func(string, string) []models.Image) {
	cm.Base64Cache.InitialImagesCache(jobName, pageNumber, getDatasetsFn, getImagesFn)
}

func (cm *CacheManager) GetImagesCache(jobName string, pageIndex int, pageNumber int, getDatasetsFn func(string) []string, getImagesFn func(string, string) []models.Image) (models.ImageItems, bool) {
	return cm.Base64Cache.GetImagesCache(jobName, pageIndex, pageNumber, getDatasetsFn, getImagesFn)
}

func (cm *CacheManager) ClearImagesCache(jobName string) {
	cm.Base64Cache.ClearImagesCache(jobName)
}

func (cm *CacheManager) ImageCacheJobExists(jobName string) bool {
	return cm.Base64Cache.ImageCacheJobExists(jobName)
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

// 新增：取得指定 job 的最大頁數
func (cm *CacheManager) GetJobMaxPages(jobName string) int {
	pageDetail := cm.Base64Cache.GetAllPageDetail(jobName)
	if pageDetail == nil {
		return 0
	}
	return len(pageDetail)
}
