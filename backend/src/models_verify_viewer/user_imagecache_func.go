package models_verify_viewer

import (
	"log"

	"github.com/patrickmn/go-cache"
)

func (cm *CacheManager) SetImageCacheStore(jobName string) {
	cacheData := NewBase64ImageCache(jobName)
	cm.ImageCacheStore.Set(jobName, cacheData, cache.DefaultExpiration)
}

func (cm *CacheManager) UpdateImageCacheStore(jobName string, cacheData *Base64ImageCache) {
	cm.ImageCacheStore.Set(jobName, cacheData, cache.DefaultExpiration)
}

func (cm *CacheManager) GetImageCacheStore(jobName string) (*Base64ImageCache, bool) {
	data, found := cm.ImageCacheStore.Get(jobName)
	if !found {
		return nil, false
	}

	if typedData, ok := data.(*Base64ImageCache); ok {
		return typedData, true
	}

	return nil, false
}

func (cm *CacheManager) ExistsImageCacheStore(jobName string) bool {
	_, found := cm.ImageCacheStore.Get(jobName)
	return found
}

func (cm *CacheManager) ClearImageCacheStore(jobName string) {
	cm.ImageCacheStore.Delete(jobName)
	log.Printf("Cleared cache for job: %s", jobName)
}

func (cm *CacheManager) CheckImageCacheStoreStats() {
	itemCount := cm.ImageCacheStore.ItemCount()
	log.Printf("[Memory] Total cached jobs: %d", itemCount)

	// Log details for each cached job
	items := cm.ImageCacheStore.Items()
	for jobName, item := range items {
		if cacheData, ok := item.Object.(*Base64ImageCache); ok {
			log.Printf("[Memory] Job: %s - Cached images: %d/%d",
				jobName, cacheData.Len(), cacheData.maxImages)
		}
	}
}

// Review cache management functions
func (cm *CacheManager) SetReviewImageCacheStore(cacheKey string) {
	cacheData := NewBase64ImageCacheWithLimit(cacheKey, maxImagesPerJobReview)
	cm.ReviewCacheStore.Set(cacheKey, cacheData, cache.DefaultExpiration)
}

func (cm *CacheManager) GetReviewImageCacheStore(cacheKey string) (*Base64ImageCache, bool) {
	data, found := cm.ReviewCacheStore.Get(cacheKey)
	if !found {
		return nil, false
	}

	if typedData, ok := data.(*Base64ImageCache); ok {
		return typedData, true
	}

	return nil, false
}

func (cm *CacheManager) ExistsReviewImageCacheStore(cacheKey string) bool {
	_, found := cm.ReviewCacheStore.Get(cacheKey)
	return found
}

func (cm *CacheManager) ClearReviewImageCacheStore(cacheKey string) {
	cm.ReviewCacheStore.Delete(cacheKey)
	log.Printf("Cleared review cache for key: %s", cacheKey)
}

func (cm *CacheManager) CheckReviewCacheStoreStats() {
	itemCount := cm.ReviewCacheStore.ItemCount()
	log.Printf("[Memory] Total review cache keys: %d", itemCount)

	// Log details for each review cache
	items := cm.ReviewCacheStore.Items()
	for cacheKey, item := range items {
		if cacheData, ok := item.Object.(*Base64ImageCache); ok {
			log.Printf("[Memory] Review Cache: %s - Cached images: %d/%d",
				cacheKey, cacheData.Len(), cacheData.maxImages)
		}
	}
}

// JobName returns the job name
func (cache *Base64ImageCache) JobName() string {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.jobName
}

// Len returns the number of cached images
func (cache *Base64ImageCache) Len() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.imageMap)
}

// Set adds or updates a single image in the cache
func (cache *Base64ImageCache) Set(imageName, base64Image string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if _, exists := cache.imageMap[imageName]; !exists {
		cache.imageMap[imageName] = base64Image
		cache.updateAccessOrder(imageName)
	}
}

// SetBatch adds multiple images to the cache
func (cache *Base64ImageCache) SetBatch(imagePaths []string, base64Images []string) {
	if len(imagePaths) != len(base64Images) {
		log.Println("Error: Image path set and Base64 image set lengths do not match")
		return
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	for idx, imagePath := range imagePaths {
		if _, exists := cache.imageMap[imagePath]; !exists {
			// Only cache non-empty base64 strings
			if base64Images[idx] != "" {
				cache.imageMap[imagePath] = base64Images[idx]
				cache.updateAccessOrder(imagePath)
			}
		}
	}

	// Clean up old images if cache is too large
	cache.cleanupOldImagesIfNeeded()
}

// Get retrieves a single image from the cache
func (cache *Base64ImageCache) Get(imagePath string) (string, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if base64, exists := cache.imageMap[imagePath]; exists {
		// Update access order for LRU
		cache.updateAccessOrder(imagePath)
		return base64, true
	}

	return "", false
}

// GetBatch retrieves multiple images from the cache
func (cache *Base64ImageCache) GetBatch(imagePaths []string) []string {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	base64Images := make([]string, 0, len(imagePaths))
	for _, imagePath := range imagePaths {
		if base64, exists := cache.imageMap[imagePath]; exists {
			base64Images = append(base64Images, base64)
			// Update access order for LRU
			cache.updateAccessOrder(imagePath)
		} else {
			base64Images = append(base64Images, "")
		}
	}

	return base64Images
}

// RemoveByPaths removes specific images from the cache by their paths
func (cache *Base64ImageCache) RemoveByPaths(imagePaths []string) int {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	removedCount := 0
	for _, imagePath := range imagePaths {
		if _, exists := cache.imageMap[imagePath]; exists {
			delete(cache.imageMap, imagePath)
			removedCount++
		}
	}

	if removedCount > 0 {
		log.Printf("Removed %d images from cache for job %s", removedCount, cache.jobName)
	}

	return removedCount
}

// CleanupEmpty removes empty/invalid entries from cache
func (cache *Base64ImageCache) CleanupEmpty() int {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	removedCount := 0
	newAccessOrder := make([]string, 0)

	for _, imagePath := range cache.accessOrder {
		if base64, exists := cache.imageMap[imagePath]; exists {
			if base64 == "" {
				// Remove empty entries
				delete(cache.imageMap, imagePath)
				removedCount++
			} else {
				newAccessOrder = append(newAccessOrder, imagePath)
			}
		}
	}

	cache.accessOrder = newAccessOrder
	return removedCount
}

// updateAccessOrder updates the LRU access order for an image (must be called with lock held)
func (cache *Base64ImageCache) updateAccessOrder(imagePath string) {
	// Remove from current position if exists
	for i, path := range cache.accessOrder {
		if path == imagePath {
			cache.accessOrder = append(cache.accessOrder[:i], cache.accessOrder[i+1:]...)
			break
		}
	}
	// Add to end (most recently used)
	cache.accessOrder = append(cache.accessOrder, imagePath)
}

// cleanupOldImagesIfNeeded removes old images when cache exceeds max size (must be called with lock held)
func (cache *Base64ImageCache) cleanupOldImagesIfNeeded() {
	if len(cache.imageMap) <= cache.maxImages {
		return
	}

	// Calculate how many images to remove
	toRemove := len(cache.imageMap) - cache.maxImages

	// Remove oldest images (from the beginning of accessOrder)
	for i := 0; i < toRemove && i < len(cache.accessOrder); i++ {
		oldPath := cache.accessOrder[i]
		delete(cache.imageMap, oldPath)
	}

	// Update access order
	cache.accessOrder = cache.accessOrder[toRemove:]
}
