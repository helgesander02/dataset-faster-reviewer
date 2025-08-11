package models_verify_viewer

import (
	"log"

	"github.com/patrickmn/go-cache"
)

// go-cache function
func (cm *CacheManager) SetImageCacheStore(jobName string) {
	cacheData := NewBase64ImageCache()
	cacheData.FillJobName(jobName)

	cm.ImageCacheStore.Set(jobName, cacheData, cache.DefaultExpiration)
}

func (cm *CacheManager) UpdateImageCacheStore(jobName string, cacheData Base64ImageCache) {
	cm.ImageCacheStore.Set(jobName, cacheData, cache.DefaultExpiration)
}

func (cm *CacheManager) GetImageCacheStore(jobName string) (Base64ImageCache, bool) {
	if data, found := cm.ImageCacheStore.Get(jobName); found {
		if typedData, ok := data.(Base64ImageCache); ok {
			return typedData, true
		}
	}
	return Base64ImageCache{}, false
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
	log.Printf("Total cached jobs: %d", itemCount)
}

// base64 image cache functions
func (base64image_cache *Base64ImageCache) FillJobName(new_job_name string) {
	base64image_cache.JobName = new_job_name
}

func (base64image_cache *Base64ImageCache) FillBase64ImageMap(new_image_name, new_base64image_name string) {
	_, exists := base64image_cache.Base64ImageMap[new_image_name]
	if !exists {
		base64image_cache.Base64ImageMap[new_image_name] = new_base64image_name
		log.Printf("New image (%s) added to Base64ImageMap (%s).\n", new_image_name, new_base64image_name)
	} else {
		log.Println("Image name already exists in Base64ImageMap. Skipping insertion.")
	}
}

func (base64image_cache *Base64ImageCache) SetBase64ImageCacheByImagePathSet(current_page_imagepath_set []string, current_page_base64image_set []string) {
	if len(current_page_imagepath_set) != len(current_page_base64image_set) {
		log.Println("Error: Image path set and Base64 image set lengths do not match.")
		return
	}

	for idx, imagePath := range current_page_imagepath_set {
		base64image_cache.FillBase64ImageMap(imagePath, current_page_base64image_set[idx])
	}
}

func (base64image_cache Base64ImageCache) GetBase64ImageCacheByImagePathSet(current_page_imagepath_set []string) []string {
	var current_page_base64image_set []string
	for _, imagePath := range current_page_imagepath_set {
		current_page_base64image_set = append(current_page_base64image_set, base64image_cache.Base64ImageMap[imagePath])
	}

	return current_page_base64image_set
}

func (base64image_cache Base64ImageCache) GetBase64ImageCacheByImagePath(current_page_imagepath string) string {
	return base64image_cache.Base64ImageMap[current_page_imagepath]
}
