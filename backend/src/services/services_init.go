package services

import (
	"backend/src/cache"
	"backend/src/models"
	"backend/src/utils"
	"log"
)

type DataManager struct {
	ImageRoot         string
	ParentData        models.Parent
	CacheManager      *cache.CacheManager
	PendingReviewData models.PendingReview
}

func NewDataManager(root string) *DataManager {
	dm := &DataManager{
		ImageRoot:         root,
		ParentData:        models.NewParentData(),
		PendingReviewData: models.NewPendingReview(),
	}
	
	// 建立 ImageProcessor
	imageProcessor := utils.NewImageProcessor()
	
	// 建立 CacheManager，傳入 DataProvider 和 ImageProcessor
	dm.CacheManager = cache.NewCacheManager(dm, imageProcessor)
	
	return dm
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
		imageProcessor := utils.NewImageProcessor()
		dm.CacheManager = cache.NewCacheManager(dm, imageProcessor)
		log.Println("Initialized CacheManager")
	}
	if dm.PendingReviewData.Items == nil {
		dm.PendingReviewData = models.NewPendingReview()
		log.Println("Initialized PendingReviewData")
	}
}
