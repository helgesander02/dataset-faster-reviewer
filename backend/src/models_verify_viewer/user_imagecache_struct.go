package models_verify_viewer

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultCacheExpiration = 10 * time.Minute // Cache expiration time
	defaultCleanupInterval = 3 * time.Minute  // Cleanup interval for expired caches
	maxImagesPerJob        = 2000             // Maximum number of images to cache per job
	maxImagesPerJobReview  = 1000             // Maximum number of images to cache for review modal
)

type CacheManager struct {
	ImageCacheStore  *cache.Cache
	ReviewCacheStore *cache.Cache
}

type Base64ImageCache struct {
	jobName     string
	imageMap    map[string]string
	accessOrder []string
	maxImages   int
	mu          sync.RWMutex
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		ImageCacheStore:  cache.New(defaultCacheExpiration, defaultCleanupInterval),
		ReviewCacheStore: cache.New(defaultCacheExpiration, defaultCleanupInterval),
	}
}

func NewBase64ImageCache(jobName string) *Base64ImageCache {
	return &Base64ImageCache{
		jobName:     jobName,
		imageMap:    make(map[string]string),
		accessOrder: make([]string, 0),
		maxImages:   maxImagesPerJob,
	}
}

func NewBase64ImageCacheWithLimit(jobName string, maxImages int) *Base64ImageCache {
	return &Base64ImageCache{
		jobName:     jobName,
		imageMap:    make(map[string]string),
		accessOrder: make([]string, 0),
		maxImages:   maxImages,
	}
}
