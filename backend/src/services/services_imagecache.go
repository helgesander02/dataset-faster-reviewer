package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"fmt"
	"log"
	"strings"
	"sync"
)

func (us *UserServices) SetBase64ImageCache(jobName string) {
	if !us.CacheManager.ExistsImageCacheStore(jobName) {
		us.CacheManager.SetImageCacheStore(jobName)
		log.Println("Image cache store created for job:", jobName)
	} else {
		log.Println("Image cache store already exists for job:", jobName)
	}
}

func (us *UserServices) GetBase64ImageCacheByPage(jobName string, pageIndex int) ([]string, []string) {
	log.Printf("[GetBase64ImageCacheByPage] START - job: %s, page: %d", jobName, pageIndex)

	imageCache := us.getOrCreateImageCache(jobName)
	if imageCache == nil {
		log.Printf("[GetBase64ImageCacheByPage] ERROR - imageCache is nil for job: %s", jobName)
		return nil, nil
	}

	// Periodic cleanup of empty entries
	imageCache.CleanupEmpty()

	imagePaths := us.CurrentPageData.ImagePathsAt(pageIndex)
	log.Printf("[GetBase64ImageCacheByPage] Retrieved %d image paths from CurrentPageData for job: %s, page: %d", len(imagePaths), jobName, pageIndex)
	if len(imagePaths) > 0 {
		log.Printf("[GetBase64ImageCacheByPage] First image path: %s", imagePaths[0])
	}

	base64Images := imageCache.GetBatch(imagePaths)
	log.Printf("[GetBase64ImageCacheByPage] Retrieved %d base64 images from cache", len(base64Images))

	// Check for empty base64 strings
	emptyCount := 0
	for i, b64 := range base64Images {
		if b64 == "" {
			emptyCount++
			if i < 3 { // Log first 3 empty entries
				log.Printf("[GetBase64ImageCacheByPage] WARNING: Empty base64 at index %d for path: %s", i, imagePaths[i])
			}
		}
	}
	if emptyCount > 0 {
		log.Printf("[GetBase64ImageCacheByPage] WARNING: Found %d empty base64 strings out of %d total", emptyCount, len(base64Images))
	}

	return imagePaths, base64Images
}

func (us *UserServices) getOrCreateImageCache(jobName string) *models_verify_viewer.Base64ImageCache {
	imageCache, exist := us.CacheManager.GetImageCacheStore(jobName)
	if exist {
		return imageCache
	}

	us.CacheManager.SetImageCacheStore(jobName)
	imageCache, exist = us.CacheManager.GetImageCacheStore(jobName)
	if !exist {
		log.Println("Failed to create image cache store for job:", jobName)
		return nil
	}

	return imageCache
}

func (us *UserServices) GetBase64ImageByPath(jobName string, imagePath string) string {
	imageCache, exist := us.CacheManager.GetImageCacheStore(jobName)
	if !exist {
		log.Println("Image cache store not found for job:", jobName)
		return ""
	}

	base64, _ := imageCache.Get(imagePath)
	return base64
}

func (us *UserServices) SetBase64ImageCacheByPage(jobName string, pageIndex int) ([]string, []string) {
	// Create unique key for this page
	lockKey := fmt.Sprintf("%s::%d", jobName, pageIndex)

	// Get or create mutex for this specific page
	lockInterface, _ := us.processingLocks.LoadOrStore(lockKey, &sync.Mutex{})
	lock := lockInterface.(*sync.Mutex)

	lock.Lock()
	defer lock.Unlock()
	defer us.processingLocks.Delete(lockKey) // Clean up after processing

	// Double-check cache after acquiring lock (another request may have completed)
	if us.ImageCacheExists(jobName, pageIndex) {
		log.Printf("[Cache] Images already cached while waiting for lock: %s, page: %d", jobName, pageIndex)
		return us.GetBase64ImageCacheByPage(jobName, pageIndex)
	}

	us.ensureImageCacheExists(jobName)

	log.Printf("[SetBase64ImageCacheByPage] START - job: %s, page: %d", jobName, pageIndex)

	// CRITICAL: Verify CurrentPageData matches the requested job
	currentJobName := us.CurrentPageData.JobName()
	if currentJobName != jobName {
		log.Printf("[SetBase64ImageCacheByPage] ERROR - Job mismatch! Requested: %s, CurrentPageData: %s. Job was switched during processing.", jobName, currentJobName)
		return nil, nil // Return nil to indicate error
	}

	imagePaths := us.CurrentPageData.ImagePathsAt(pageIndex)
	log.Printf("[SetBase64ImageCacheByPage] Retrieved %d image paths for job: %s, page: %d", len(imagePaths), jobName, pageIndex)
	if len(imagePaths) > 0 {
		log.Printf("[SetBase64ImageCacheByPage] First image path: %s", imagePaths[0])

		// Double-check: verify first path contains job name as additional safety
		if !strings.Contains(imagePaths[0], jobName) {
			log.Printf("[SetBase64ImageCacheByPage] ERROR - Path validation failed! Path '%s' doesn't contain job name '%s'", imagePaths[0], jobName)
			return nil, nil
		}
	}

	base64Images := utils.CompressImageSetToBase64(imagePaths, pageIndex)
	log.Printf("[SetBase64ImageCacheByPage] Compressed %d images, result count: %d", len(imagePaths), len(base64Images))

	// Check if processing was cancelled or failed (all empty results)
	if isResultSetEmpty(base64Images) {
		log.Printf("[Cache] Task was cancelled or failed for job: %s, page: %d - returning nil to trigger retry", jobName, pageIndex)
		return nil, nil // Return nil to signal handler that request should fail/retry
	}

	// Check for individual empty base64 strings
	emptyCount := 0
	for i, b64 := range base64Images {
		if b64 == "" {
			emptyCount++
			if i < 3 { // Log first 3 empty entries
				log.Printf("[SetBase64ImageCacheByPage] WARNING: Empty base64 at index %d for path: %s", i, imagePaths[i])
			}
		}
	}
	if emptyCount > 0 {
		log.Printf("[SetBase64ImageCacheByPage] WARNING: Found %d empty base64 strings out of %d total for job: %s, page: %d", emptyCount, len(base64Images), jobName, pageIndex)
	}

	cacheData, found := us.CacheManager.GetImageCacheStore(jobName)
	if !found {
		log.Println("Image cache store not found for job:", jobName)
		return nil, nil
	}

	cacheData.SetBatch(imagePaths, base64Images)
	log.Printf("[Cache] Successfully cached %d images for job: %s, page: %d (with %d empty strings)", len(imagePaths), jobName, pageIndex, emptyCount)
	return imagePaths, base64Images
}

// isResultSetEmpty checks if all images in the result are empty (cancelled/failed)
func isResultSetEmpty(base64Images []string) bool {
	if len(base64Images) == 0 {
		return true
	}

	// Check if all images are empty
	for _, img := range base64Images {
		if img != "" {
			return false
		}
	}

	return true
}

func (us *UserServices) ensureImageCacheExists(jobName string) {
	if !us.CacheManager.ExistsImageCacheStore(jobName) {
		us.CacheManager.SetImageCacheStore(jobName)
	}
}

func (us *UserServices) ImageCacheExists(jobName string, pageIndex int) bool {
	data, found := us.CacheManager.GetImageCacheStore(jobName)
	if !found {
		log.Println("Image cache store not found for job:", jobName)
		return false
	}

	return us.allImagesAreCached(data, jobName, pageIndex)
}

func (us *UserServices) allImagesAreCached(data *models_verify_viewer.Base64ImageCache, jobName string, pageIndex int) bool {
	_, imagePaths := us.GetImageCacheByPage(jobName, pageIndex)
	for _, imagePath := range imagePaths {
		if !us.isImageCached(data, imagePath, jobName) {
			return false
		}
	}
	return true
}

func (us *UserServices) isImageCached(data *models_verify_viewer.Base64ImageCache, imagePath, jobName string) bool {
	base64Image, found := data.Get(imagePath)
	if !found || base64Image == "" {
		log.Printf("Image path %s not found in cache for job %s", imagePath, jobName)
		return false
	}
	return true
}

// GetOriginalImageBase64 returns the original, uncompressed image in base64 format
func (us *UserServices) GetOriginalImageBase64(imagePath string) (string, error) {
	return utils.ImageToBase64(imagePath)
}

// RemoveImagesFromCache removes deleted images from the image cache for a specific job
func (us *UserServices) RemoveImagesFromCache(jobName string, imagePaths []string) int {
	if len(imagePaths) == 0 {
		return 0
	}

	imageCache, found := us.CacheManager.GetImageCacheStore(jobName)
	if !found {
		log.Printf("Image cache not found for job %s, skipping cache cleanup", jobName)
		return 0
	}

	removedCount := imageCache.RemoveByPaths(imagePaths)
	return removedCount
}

func (us *UserServices) ClearImageCache(job string) {
	us.CacheManager.ClearImageCacheStore(job)
}
