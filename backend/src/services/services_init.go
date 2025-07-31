package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
)

//type JointInformation struct {
//	ImageRoot string
//}

//type JointServices struct {
//	JobList models_verify_viewer.JobList
//}

//type UserServices struct {
//	CacheManager      *cache.CacheManager
//	CurrentPageData   models_verify_viewer.ImagesPerPageCache
//	PendingReviewData models_verify_viewer.PendingReview
//}

type DataManager struct {
	ImageRoot         string
	JobList           models_verify_viewer.JobList
	CacheManager      *models_verify_viewer.CacheManager
	PendingReviewData models_verify_viewer.PendingReview
}

func NewDataManager(root string) *DataManager {
	dm := &DataManager{
		ImageRoot:         root,
		JobList:           models_verify_viewer.NewJobList(),
		PendingReviewData: models_verify_viewer.NewPendingReview(),
	}

	// 建立 ImageProcessor
	imageProcessor := utils.NewImageProcessor()

	// 建立 CacheManager，傳入 DataProvider 和 ImageProcessor
	dm.CacheManager = models_verify_viewer.NewCacheManager(dm, imageProcessor)

	return dm
}

func (dm *DataManager) SetupServices() {
	if dm.ImageRoot == "" {
		log.Println("Image root is not set")
	}
	if dm.JobList.Jobs == nil {
		dm.JobList = models_verify_viewer.NewJobList()
		log.Println("Initialized ParentData")
	}
	if dm.CacheManager == nil {
		imageProcessor := utils.NewImageProcessor()
		dm.CacheManager = models_verify_viewer.NewCacheManager(dm, imageProcessor)
		log.Println("Initialized CacheManager")
	}
	if dm.PendingReviewData.Items == nil {
		dm.PendingReviewData = models_verify_viewer.NewPendingReview()
		log.Println("Initialized PendingReviewData")
	}
}
