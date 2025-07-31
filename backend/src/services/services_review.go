package services

import (
	"backend/src/models_verify_viewer"
	"log"
)

func (dm *DataManager) SavePendingReviewData(body interface{}) int {
	itemsData, ok := body.([]interface{})
	if !ok {
		log.Println("SavePendingReview: invalid data format")
		return 0
	}

	var items []models_verify_viewer.PendingReviewItem
	for _, item := range itemsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		pendingItem := models_verify_viewer.PendingReviewItem{
			Job:       getString(itemMap, "job"),
			Dataset:   getString(itemMap, "dataset"),
			ImageName: getString(itemMap, "imageName"),
			ImagePath: getString(itemMap, "imagePath"),
		}
		items = append(items, pendingItem)
	}

	pending := models_verify_viewer.NewPendingReview()
	if len(items) == 0 {
		log.Println("Empty list provided")
		dm.PendingReviewData.MergePendingReviewItems(pending)
		log.Println("SavePendingReview: cleaned up pending review items")
		return -1
	}

	pending.Items = items
	dm.PendingReviewData.MergePendingReviewItems(pending)
	log.Printf("SavePendingReview: loaded %d items", len(items))
	return len(items)
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (dm *DataManager) GetPendingReviewItems() []models_verify_viewer.PendingReviewItem {
	return dm.PendingReviewData.Items
}

func (dm *DataManager) ClearPendingReviewData() {
	dm.PendingReviewData.ClearPendingReviewItems()
}
