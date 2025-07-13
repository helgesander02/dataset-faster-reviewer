package services

import (
	"backend/models"
	"log"
)

type DataManager struct {
	ImageRoot          string
	ParentData         models.Parent
	JobCache           models.JobCache
	ImagesPerPageCache models.ImagesPerPageCache
	PendingReviewData  models.PendingReview
}

func NewDataManager(root string) *DataManager {
	return &DataManager{
		ImageRoot:          root,
		ParentData:         models.NewParentData(),
		JobCache:           models.NewJobCache(),
		ImagesPerPageCache: models.NewImagesPerPageCache(),
		PendingReviewData:  models.NewPendingReview(),
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
	if dm.JobCache.Jobs == nil {
		dm.JobCache = models.NewJobCache()
		log.Println("Initialized JobCache")
	}
	if dm.ImagesPerPageCache.Pages == nil {
		dm.ImagesPerPageCache = models.NewImagesPerPageCache()
		log.Println("Initialized ImagesCache")
	}
	if dm.PendingReviewData.Items == nil {
		dm.PendingReviewData = models.NewPendingReview()
		log.Println("Initialized PendingReviewData")
	}
}
