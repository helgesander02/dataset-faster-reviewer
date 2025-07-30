package cache

import (
	"backend/src/models"
	"log"
)

type CacheManager struct {
	JobCache    *JobCache
	Base64Cache *Base64Cache
}

func NewCacheManager(dataProvider DataProvider, imageProcessor ImageProcessor) *CacheManager {
	return &CacheManager{
		JobCache:    NewJobCache(3),
		Base64Cache: NewBase64Cache(dataProvider, imageProcessor),
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
func (cm *CacheManager) GetAllPageDetail(jobName string) (map[int]string, error) {
	return cm.Base64Cache.GetAllPageDetail(jobName)
}

func (cm *CacheManager) InitializeImagesCache(jobName string, pageSize int) error {
	return cm.Base64Cache.InitializeImagesCache(jobName, pageSize)
}

func (cm *CacheManager) GetImagesCache(jobName string, pageIndex int) (models.ImageItems, error) {
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
