package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
)

const (
	reviewCacheKey       = "pending_review"
	maxReviewCacheImages = 1000
)

// GetOrCreateReviewBase64Images retrieves or creates base64 images for review modal.
// Uses a separate cache (ReviewCacheStore) with 1000 image limit to avoid conflicts with ImageGrid cache.
func (us *UserServices) GetOrCreateReviewBase64Images(imagePaths []string) []string {
	if len(imagePaths) == 0 {
		return []string{}
	}

	// Get or create review cache
	reviewCache := us.getOrCreateReviewCache()
	if reviewCache == nil {
		log.Println("ERROR: Failed to create review cache")
		return make([]string, len(imagePaths))
	}

	// Check cache for existing images
	base64Images := make([]string, len(imagePaths))
	missingPaths := make([]string, 0)
	missingIndices := make([]int, 0)

	for i, imagePath := range imagePaths {
		if base64, found := reviewCache.Get(imagePath); found && base64 != "" {
			base64Images[i] = base64
		} else {
			missingPaths = append(missingPaths, imagePath)
			missingIndices = append(missingIndices, i)
		}
	}

	// Compress missing images
	if len(missingPaths) > 0 {
		compressedImages := us.compressReviewImages(missingPaths)

		// Store in cache and populate result
		for i, compressedImage := range compressedImages {
			if compressedImage != "" {
				originalIndex := missingIndices[i]
				imagePath := missingPaths[i]

				base64Images[originalIndex] = compressedImage
				reviewCache.Set(imagePath, compressedImage)
			}
		}

		// Update cache in CacheManager
		us.CacheManager.ReviewCacheStore.Set(reviewCacheKey, reviewCache, 0)
	}

	return base64Images
}

// getOrCreateReviewCache retrieves or creates the review cache
func (us *UserServices) getOrCreateReviewCache() *models_verify_viewer.Base64ImageCache {
	reviewCache, exist := us.CacheManager.GetReviewImageCacheStore(reviewCacheKey)
	if exist {
		return reviewCache
	}

	us.CacheManager.SetReviewImageCacheStore(reviewCacheKey)
	reviewCache, exist = us.CacheManager.GetReviewImageCacheStore(reviewCacheKey)
	if !exist {
		log.Println("Failed to create review cache store")
		return nil
	}

	return reviewCache
}

// compressReviewImages compresses images for review modal (synchronous compression)
func (us *UserServices) compressReviewImages(imagePaths []string) []string {
	// Use synchronous compression for review images to avoid task manager conflicts
	base64Images := make([]string, len(imagePaths))

	for i, imagePath := range imagePaths {
		base64Image, err := utils.CompressImageToBase64(imagePath)
		if err != nil {
			log.Printf("Failed to compress image %s: %v", imagePath, err)
			base64Images[i] = ""
		} else {
			base64Images[i] = base64Image
		}
	}

	return base64Images
}
