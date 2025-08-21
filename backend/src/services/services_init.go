package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
)

var (
	ImageRoot string
	BackupDir string
)

type JointServices struct {
	JobList           models_verify_viewer.JobList
	PendingReviewData models_verify_viewer.PendingReview
}

func NewJointServices() *JointServices {
	JointServices := &JointServices{
		JobList:           models_verify_viewer.NewJobList(),
		PendingReviewData: models_verify_viewer.NewPendingReview(BackupDir),
	}

	utils.ConcurrentJobScanner(ImageRoot, &JointServices.JobList.Jobs)
	return JointServices
}

type UserServices struct {
	CacheManager    *models_verify_viewer.CacheManager
	CurrentPageData models_verify_viewer.Pages
}

func NewUserServices() *UserServices {
	us := &UserServices{
		CacheManager:    models_verify_viewer.NewCacheManager(),
		CurrentPageData: models_verify_viewer.NewPages(),
	}

	return us
}

func CheckServicesState(us *UserServices, js *JointServices) {
	if ImageRoot == "" {
		log.Println("Image root is not set")

	} else if js.JobList.Jobs == nil {
		log.Println("Initialized ParentData")

	} else if js.PendingReviewData.Items == nil {
		log.Println("Initialized PendingReviewData")

	} else if us.CacheManager == nil {
		log.Println("Initialized CacheManager")

	} else if us.CurrentPageData.PageItems == nil {
		log.Println("Initialized CurrentPageData")

	} else {
		log.Println("Anything is initialized")
	}
}
