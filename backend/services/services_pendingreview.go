package services


import (
	"log"
    "encoding/json"
	"backend/models"
)


func (dm *DataManager) SavePendingReviewData(body *json.Decoder) int {
    items := models.NewPendingReviewItems()
    if err := body.Decode(&items); err != nil {
        log.Printf("SavePendingReview: decode error: %v", err)
        log.Println("Invalid request format")
        return 0 
    }

    pending := models.NewPendingReview()
    if len(items) == 0 {
        log.Println("Empty list provided")
        dm.PendingReviewData.MergePendingReviewItems(pending)
        log.Println("SavePendingReview: cleaned up pending review items")
        return -1

    } else {
        pending.Items = items
        dm.PendingReviewData.MergePendingReviewItems(pending)
        log.Println(dm.PendingReviewData)
        log.Printf("SavePendingReview: loaded %d items across %d jobs", len(items), len(pending.Items))
        return len(items) 

    }
}


func (dm *DataManager) GetPendingReviewItems() []models.PendingReviewItem {
	return dm.PendingReviewData.Items
}


func (dm *DataManager) ClearPendingReviewData() {
	dm.PendingReviewData.ClearPendingReviewItems()
}
