package services

import (
	"backend/src/models_verify_viewer"
	"log"
	"path/filepath"
)

func (js *JointServices) SavePendingReviewData(body interface{}) int {
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

	pending := models_verify_viewer.NewPendingReview(BackupDir)
	if len(items) == 0 {
		log.Println("Empty list provided")
		js.PendingReviewData.MergePendingReviewItems(pending)
		log.Println("SavePendingReview: cleaned up pending review items")
		return -1
	}

	pending.Items = items
	js.PendingReviewData.MergePendingReviewItems(pending)
	if err := js.PendingReviewData.CreateBackup(); err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	}
	log.Printf("SavePendingReview: loaded %d items", len(items))
	return len(items)
}

func (js *JointServices) GetBackupList() ([]models_verify_viewer.BackupInfo, error) {
	return js.PendingReviewData.ListBackups()
}

func (js *JointServices) RestoreFromBackup(filename string) error {
	restoredData, err := js.PendingReviewData.RestoreFromBackup(filename)
	if err != nil {
		return err
	}

	js.PendingReviewData = restoredData
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

func (js *JointServices) GetPendingReviewItems() []models_verify_viewer.PendingReviewItem {
	return js.PendingReviewData.Items
}

func (js *JointServices) ClearPendingReviewData() {
	if err := js.PendingReviewData.CreateBackup(); err != nil {
		log.Printf("Warning: Failed to create backup before clear: %v", err)
	}

	js.PendingReviewData.ClearPendingReviewItems()
}

func (js *JointServices) GetPendingReviewImagePaths() []string {
	items := js.GetPendingReviewItems()
	imagePaths := make([]string, 0, len(items))

	for _, item := range items {
		fullPath := filepath.Join(ImageRoot, item.JobName, item.DatasetName, item.ImageName)
		imagePaths = append(imagePaths, fullPath)
	}
	log.Printf("GetPendingReviewImagePaths: found %d image paths", len(imagePaths))

	return imagePaths
}
