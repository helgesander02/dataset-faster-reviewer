package services

import (
	"backend/src/models_verify_viewer"
	"log"
	"path/filepath"
)

func (us *UserServices) SavePendingReviewData(body interface{}) int {
	itemsData, ok := body.([]interface{})
	if !ok {
		log.Println("SavePendingReview: invalid data format")
		return 0
	}

	items := models_verify_viewer.NewPendingReviewItemSet()
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

	pending := models_verify_viewer.NewPendingReview()
	if len(items) == 0 {
		log.Println("Empty list provided")
		us.PendingReviewData.MergePendingReviewItems(pending)
		log.Println("SavePendingReview: cleaned up pending review items")
		return -1
	}

	pending.Items = items
	us.PendingReviewData.MergePendingReviewItems(pending)
	if err := us.BackupManager.CreateBackup(us.PendingReviewData); err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	}
	log.Printf("SavePendingReview: loaded %d items", len(items))
	return len(items)
}

func (us *UserServices) GetBackupList() ([]models_verify_viewer.BackupInfo, error) {
	return us.BackupManager.ListBackups()
}

func (us *UserServices) RestoreFromBackup(filename string) error {
	restoredData, err := us.BackupManager.RestoreFromBackup(filename)
	if err != nil {
		return err
	}

	if err := us.BackupManager.CreateBackup(us.PendingReviewData); err != nil {
		log.Printf("Warning: Failed to create backup before restore: %v", err)
	}

	us.PendingReviewData = restoredData
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

func (us *UserServices) GetPendingReviewItems() []models_verify_viewer.PendingReviewItem {
	return us.PendingReviewData.Items
}

func (us *UserServices) ClearPendingReviewData() {
	if err := us.BackupManager.CreateBackup(us.PendingReviewData); err != nil {
		log.Printf("Warning: Failed to create backup before clear: %v", err)
	}

	us.PendingReviewData.ClearPendingReviewItems()
}

func (us *UserServices) GetPendingReviewImagePaths() []string {
	items := us.GetPendingReviewItems()
	imagePaths := make([]string, 0, len(items))

	for _, item := range items {
		fullPath := filepath.Join(ImageRoot, item.JobName, item.DatasetName, item.ImageName)
		imagePaths = append(imagePaths, fullPath)
	}
	log.Printf("GetPendingReviewImagePaths: found %d image paths", len(imagePaths))

	return imagePaths
}
