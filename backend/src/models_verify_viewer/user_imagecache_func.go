// user_imagecache_func.go
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

// Page Cache 相關方法 - 委託給 PageCache
func (cm *CacheManager) GetAllPageDetail(jobName string) (map[int]string, error) {
	return cm.PageCache.GetAllPageDetail(jobName)
}

func (cm *CacheManager) InitializeImagesCache(jobName string, pageSize int) error {
	return cm.PageCache.InitializeImagesCache(jobName, pageSize)
}

func (cm *CacheManager) GetImagesCache(jobName string, pageIndex int) (ImageItems, error) {
	return cm.PageCache.GetImagesCache(jobName, pageIndex)
}

func (cm *CacheManager) ClearImagesCache(jobName string) {
	cm.PageCache.ClearImagesCache(jobName)
}

func (cm *CacheManager) ImageCacheJobExists(jobName string) bool {
	return cm.PageCache.ImageCacheJobExists(jobName)
}

func (cm *CacheManager) GetJobMaxPages(jobName string) (int, error) {
	return cm.PageCache.GetMaxPages(jobName)
}

// Base64 Cache 相關方法
func (bc *Base64Cache) ClearImagesCache(jobName string) {
	bc.cache.Delete(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
}

func (bc *Base64Cache) ImageCacheJobExists(jobName string) bool {
	_, found := bc.cache.Get(jobName)
	return found
}

// 清理所有快取
func (cm *CacheManager) ClearAllCache() {
	cm.JobCache.Clear()
	cm.Base64Cache.cache.Flush()
	log.Println("Cleared all caches")
}

// 取得快取統計資訊
func (cm *CacheManager) GetCacheStats() map[string]int {
	return map[string]int{
		"job_cache_items":    cm.JobCache.ItemCount(),
		"base64_cache_items": cm.Base64Cache.cache.ItemCount(),
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
