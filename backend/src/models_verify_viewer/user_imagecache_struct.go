package models_verify_viewer

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheManager struct {
	ImageCacheStore *cache.Cache
}

type Base64ImageCache struct {
	JobName        string            `json:"job_name"`
	Base64ImageMap map[string]string `json:"image_map"`
}

// set a default expiration time of 30 minutes and cleanup interval of 10 minutes
func NewCacheManager() *CacheManager {
	return &CacheManager{
		ImageCacheStore: cache.New(30*time.Minute, 10*time.Minute),
	}
}

func NewBase64ImageCache() Base64ImageCache {
	return Base64ImageCache{
		JobName:        "",
		Base64ImageMap: NewBase64ImageMap(),
	}
}

func NewBase64ImageMap() map[string]string {
	return make(map[string]string)
}
