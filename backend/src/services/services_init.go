package services

import (
	"backend/src/cache"
	"backend/src/models"
	"log"
)

type DataManager struct {
	ImageRoot         string
	ParentData        models.Parent
	CacheManager      *cache.CacheManager
	PendingReviewData models.PendingReview
}

func NewDataManager(root string) *DataManager {
	return &DataManager{
		ImageRoot:         root,
		ParentData:        models.NewParentData(),
		CacheManager:      cache.NewCacheManager(),
		PendingReviewData: models.NewPendingReview(),
	}
}

func (dm *DataManager) SetupServices() {
	if dm.ImageRoot == "" {
		log.Println("Image root is not set")
	}
	if dm.ParentData.Jobs == nil {
		dm.ParentData = models.NewParentData()
		log.Println("Initialized ParentData")
	}
	if dm.CacheManager == nil {
		dm.CacheManager = cache.NewCacheManager()
		log.Println("Initialized CacheManager")
	}
	if dm.PendingReviewData.Items == nil {
		dm.PendingReviewData = models.NewPendingReview()
		log.Println("Initialized PendingReviewData")
	}
}
