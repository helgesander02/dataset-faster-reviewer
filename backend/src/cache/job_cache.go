package cache

import (
	"backend/src/models"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type JobCache struct {
	cache    *cache.Cache
	maxItems int
}

func NewJobCache(maxItems int) *JobCache {
	// 設定 cache 過期時間為 1 小時，清理間隔為 30 分鐘
	c := cache.New(1*time.Hour, 30*time.Minute)
	return &JobCache{
		cache:    c,
		maxItems: maxItems,
	}
}

func (jc *JobCache) Get(jobName string) (models.Job, bool) {
	data, found := jc.cache.Get(jobName)
	if !found {
		return models.Job{}, false
	}
	
	job, ok := data.(models.Job)
	if !ok {
		log.Printf("Invalid job data type in cache for job: %s", jobName)
		jc.cache.Delete(jobName)
		return models.Job{}, false
	}
	
	log.Printf("Job found in cache: %s", jobName)
	return job, true
}

func (jc *JobCache) Set(job models.Job) {
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
