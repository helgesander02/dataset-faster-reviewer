package services

import (
	"backend/src/models_verify_viewer"
	"log"
)

func (us *UserServices) SavePendingReviewData(body interface{}) int {
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
		us.PendingReviewData.MergePendingReviewItems(pending)
		log.Println("SavePendingReview: cleaned up pending review items")
		return -1
	}

	pending.Items = items
	us.PendingReviewData.MergePendingReviewItems(pending)
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

func (us *UserServices) GetPendingReviewItems() []models_verify_viewer.PendingReviewItem {
	return us.PendingReviewData.Items
}

func (us *UserServices) ClearPendingReviewData() {
	us.PendingReviewData.ClearPendingReviewItems()
}
