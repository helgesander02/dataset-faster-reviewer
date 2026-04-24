package services

import (
	"backend/src/models_verify_viewer"
	"fmt"
	"log"
	"path/filepath"
)

// SavePendingReviewData saves pending review items from the request body
// Returns the number of items saved, or -1 if cleared
func (js *JointServices) SavePendingReviewData(body interface{}) int {
	items, err := parseReviewItems(body)
	if err != nil {
		log.Println("Failed to parse review items:", err)
		return 0
	}

	if len(items) == 0 {
		return js.clearPendingReviewItems()
	}

	return js.mergePendingReviewItems(items)
}

func parseReviewItems(body interface{}) ([]models_verify_viewer.PendingReviewItem, error) {
	itemsData, ok := body.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	items := make([]models_verify_viewer.PendingReviewItem, 0, len(itemsData))
	for _, item := range itemsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		pendingItem := models_verify_viewer.PendingReviewItem{
			JobName:     getString(itemMap, "job"),
			DatasetName: getString(itemMap, "dataset"),
			ImageName:   getString(itemMap, "imageName"),
			ImagePath:   getString(itemMap, "imagePath"),
		}
		items = append(items, pendingItem)
	}

	return items, nil
}

func (js *JointServices) clearPendingReviewItems() int {
	pending := models_verify_viewer.NewPendingReview()
	js.PendingReviewData.Merge(pending)
	return -1
}

func (js *JointServices) mergePendingReviewItems(items []models_verify_viewer.PendingReviewItem) int {
	pending := models_verify_viewer.NewPendingReview()
	pending.Replace(items)
	js.PendingReviewData.Merge(pending)

	backupDir := GetBackupDir()
	if err := js.PendingReviewData.CreateBackup(backupDir); err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	}

	return len(items)
}

// GetBackupList returns a list of all available backups
func (js *JointServices) GetBackupList() ([]models_verify_viewer.BackupInfo, error) {
	backupDir := GetBackupDir()
	return js.PendingReviewData.ListBackups(backupDir)
}

// RestoreFromBackup restores pending review data from a backup file
func (js *JointServices) RestoreFromBackup(filename string) error {
	backupDir := GetBackupDir()
	err := js.PendingReviewData.RestoreFromBackup(backupDir, filename)
	if err != nil {
		return err
	}

	return nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetPendingReviewItems returns all pending review items
func (js *JointServices) GetPendingReviewItems() []models_verify_viewer.PendingReviewItem {
	return js.PendingReviewData.Items()
}

// ClearPendingReview clears all pending review data
func (js *JointServices) ClearPendingReview() {
	// Clear the pending review data
	js.PendingReviewData.Clear()
	
	// Create a backup of the cleared state
	backupDir := GetBackupDir()
	if err := js.PendingReviewData.CreateBackup(backupDir); err != nil {
		log.Printf("Warning: Failed to create backup after clear: %v", err)
	}
}

// ClearPendingReviewData clears all pending review data after creating a backup
func (js *JointServices) ClearPendingReviewData() {
	backupDir := GetBackupDir()
	if err := js.PendingReviewData.CreateBackup(backupDir); err != nil {
		log.Printf("Warning: Failed to create backup before clear: %v", err)
	}

	js.PendingReviewData.Clear()
}

// GetPendingReviewImagePaths returns full file paths for all pending review items
func (js *JointServices) GetPendingReviewImagePaths() []string {
	items := js.GetPendingReviewItems()
	imagePaths := make([]string, 0, len(items))
	root := GetImageRoot()

	for _, item := range items {
		fullPath := filepath.Join(root, item.JobName, item.DatasetName, item.ImageName)
		imagePaths = append(imagePaths, fullPath)
	}

	return imagePaths
}

// DeleteImageResult contains information about the deletion operation
type DeleteImageResult struct {
	DeletedCount int      `json:"deleted_count"`
	CacheCleared bool     `json:"cache_cleared"`
	AffectedJobs []string `json:"affected_jobs"`
	DeletedPaths []string `json:"deleted_paths"`
}

// DeleteSelectedImages deletes physical image files and removes them from pending review
func (js *JointServices) DeleteSelectedImages(body interface{}) (*DeleteImageResult, error) {
	itemsData, ok := body.([]interface{})
	if !ok {
		log.Println("Invalid data format for DeleteSelectedImages")
		return &DeleteImageResult{DeletedCount: 0, CacheCleared: false}, nil
	}

	deletedItems, deletedPaths, affectedJobs := js.deletePhysicalFiles(itemsData)
	deletedCount := len(deletedItems)

	if deletedCount > 0 {
		js.removePendingReviewItems(deletedItems)
	}

	return &DeleteImageResult{
		DeletedCount: deletedCount,
		CacheCleared: false,
		AffectedJobs: affectedJobs,
		DeletedPaths: deletedPaths,
	}, nil
}

func (js *JointServices) deletePhysicalFiles(itemsData []interface{}) (map[string]bool, []string, []string) {
	deletedItems := make(map[string]bool)
	deletedPaths := make([]string, 0)
	affectedJobsMap := make(map[string]bool)
	root := GetImageRoot()

	for _, item := range itemsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		jobName := getString(itemMap, "job")
		datasetName := getString(itemMap, "dataset")
		imageName := getString(itemMap, "imageName")

		if !isValidImageItem(jobName, datasetName, imageName) {
			continue
		}

		fullPath := buildImagePath(root, jobName, datasetName, imageName)
		if js.deleteImageFile(fullPath) {
			key := createItemKey(jobName, datasetName, imageName)
			deletedItems[key] = true
			deletedPaths = append(deletedPaths, fullPath)
			affectedJobsMap[jobName] = true
		}
	}

	// Convert affected jobs map to slice
	affectedJobs := make([]string, 0, len(affectedJobsMap))
	for jobName := range affectedJobsMap {
		affectedJobs = append(affectedJobs, jobName)
	}

	return deletedItems, deletedPaths, affectedJobs
}

func isValidImageItem(jobName, datasetName, imageName string) bool {
	return jobName != "" && datasetName != "" && imageName != ""
}

func buildImagePath(root, jobName, datasetName, imageName string) string {
	return filepath.Join(root, jobName, datasetName, "image", imageName)
}

func (js *JointServices) deleteImageFile(fullPath string) bool {
	if err := models_verify_viewer.DeleteImageFile(fullPath); err != nil {
		log.Printf("Failed to delete image %s: %v", fullPath, err)
		return false
	}

	return true
}

func (js *JointServices) removePendingReviewItems(deletedItems map[string]bool) {
	items := js.PendingReviewData.Items()
	newItems := make([]models_verify_viewer.PendingReviewItem, 0)
	for _, item := range items {
		key := createItemKey(item.JobName, item.DatasetName, item.ImageName)
		if !deletedItems[key] {
			newItems = append(newItems, item)
		}
	}

	js.PendingReviewData.Replace(newItems)

	backupDir := GetBackupDir()
	if err := js.PendingReviewData.CreateBackup(backupDir); err != nil {
		log.Printf("Warning: Failed to create backup after deletion: %v", err)
	}

	log.Printf("Removed %d items from pending review list", len(deletedItems))
}

func createItemKey(jobName, datasetName, imageName string) string {
	return fmt.Sprintf("%s|%s|%s", jobName, datasetName, imageName)
}

// CleanupDeletedImagesFromCache removes deleted images from cache and page data
func (js *JointServices) CleanupDeletedImagesFromCache(us *UserServices, result *DeleteImageResult) {
	if result.DeletedCount == 0 {
		return
	}

	// Clean up page data
	pageDataRemoved := us.RemoveImagesFromPageData(result.DeletedPaths)

	// Clean up image cache for each affected job
	totalCacheRemoved := 0
	for _, jobName := range result.AffectedJobs {
		cacheRemoved := us.RemoveImagesFromCache(jobName, result.DeletedPaths)
		totalCacheRemoved += cacheRemoved
	}

	if totalCacheRemoved > 0 || pageDataRemoved > 0 {
		result.CacheCleared = true
		log.Printf("Cache cleanup complete: %d images removed from cache, %d from page data",
			totalCacheRemoved, pageDataRemoved)
	}
}
