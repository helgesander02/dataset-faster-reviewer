package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
)

var (
	ImageRoot string
)

type JointServices struct {
	JobList models_verify_viewer.JobList
}

func NewJointServices() *JointServices {
	return &JointServices{
		JobList: models_verify_viewer.NewJobList(),
	}
}

type UserServices struct {
	CacheManager      *models_verify_viewer.CacheManager
	CurrentPageData   models_verify_viewer.ImagesPerPageCache
	PendingReviewData models_verify_viewer.PendingReview
}

func NewUserServices() *UserServices {
	us := &UserServices{
		CurrentPageData:   models_verify_viewer.NewImagesPerPageCache(),
		PendingReviewData: models_verify_viewer.NewPendingReview(),
	}

	imageProcessor := utils.NewImageProcessor()
	us.CacheManager = models_verify_viewer.NewCacheManager(us, imageProcessor)

	return us
}

func CheckServicesStart(us *UserServices, js *JointServices) {
	if ImageRoot == "" {
		log.Println("Image root is not set")

	} else if js.JobList.Jobs == nil {
		log.Println("Initialized ParentData")

	} else if us.CacheManager == nil {
		log.Println("Initialized CacheManager")

	} else if us.CurrentPageData.Pages == nil {
		log.Println("Initialized CurrentPageData")

	} else if us.PendingReviewData.Items == nil {
		log.Println("Initialized PendingReviewData")

	} else {
		log.Println("Anything is initialized")
	}
}
