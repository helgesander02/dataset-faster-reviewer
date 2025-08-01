// user_imagecache_struct.go
package models_verify_viewer

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheManager struct {
	JobCache    *JobCache
	Base64Cache *Base64Cache
	PageCache   *PageCache
}

func NewCacheManager(dataProvider DataProvider, imageProcessor ImageProcessor) *CacheManager {
	return &CacheManager{
		JobCache:    NewJobCache(3),
		Base64Cache: NewBase64Cache(dataProvider, imageProcessor),
		PageCache:   NewPageCache(dataProvider, imageProcessor),
	}
}

type Base64Cache struct {
	cache          *cache.Cache
	dataProvider   DataProvider
	imageProcessor ImageProcessor
}

func NewBase64Cache(dataProvider DataProvider, imageProcessor ImageProcessor) *Base64Cache {
	c := cache.New(30*time.Minute, 10*time.Minute) // set a default expiration time of 30 minutes and cleanup interval of 10 minutes
	return &Base64Cache{
		cache:          c,
		dataProvider:   dataProvider,
		imageProcessor: imageProcessor,
	}
}

type JobCache struct {
	cache    *cache.Cache
	maxItems int
}

func NewJobCache(maxItems int) *JobCache {
	c := cache.New(1*time.Hour, 30*time.Minute) // set a default expiration time of 1 hour and cleanup interval of 30 minutes
	return &JobCache{
		cache:    c,
		maxItems: maxItems,
	}
}

type DataProvider interface {
	GetDatasets(jobName string) []string
	GetImages(jobName, datasetName string) []Image
}

type ImageProcessor interface {
	CompressToBase64(imagePath string) (string, error)
}
