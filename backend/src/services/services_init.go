package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"context"
	"log"
	"sync"
)

var (
	configOnce sync.Once
	imageRoot  string
	backupDir  string
)

func SetConfig(root, backup string) {
	configOnce.Do(func() {
		imageRoot = root
		backupDir = backup
		log.Printf("Configuration set - ImageRoot: %s, BackupDir: %s", root, backup)
	})
}

func GetImageRoot() string {
	return imageRoot
}

func GetBackupDir() string {
	return backupDir
}

type JointServices struct {
	JobList           *models_verify_viewer.JobList
	PendingReviewData *models_verify_viewer.PendingReview
}

func NewJointServices(ctx context.Context) *JointServices {
	js := &JointServices{
		JobList:           models_verify_viewer.NewJobList(),
		PendingReviewData: models_verify_viewer.NewPendingReview(),
	}

	imageRoot := GetImageRoot()
	utils.ConcurrentJobScanner(ctx, imageRoot, js.JobList)

	js.autoRestoreLatestBackup()
	return js
}

type UserServices struct {
	CacheManager    *models_verify_viewer.CacheManager
	CurrentPageData *models_verify_viewer.Pages
	processingLocks sync.Map // Per-page locks to prevent duplicate processing
}

func NewUserServices() *UserServices {
	return &UserServices{
		CacheManager:    models_verify_viewer.NewCacheManager(),
		CurrentPageData: models_verify_viewer.NewPages(),
	}
}

func (js *JointServices) autoRestoreLatestBackup() {
	backupDir := GetBackupDir()
	latestBackup, err := js.PendingReviewData.GetLatestBackup(backupDir)
	if err != nil {
		log.Printf("No backup found to restore on startup (this is normal for first run): %v", err)
		return
	}
	log.Printf("Found latest backup: %s, restoring...", latestBackup)

	err = js.PendingReviewData.RestoreFromBackup(backupDir, latestBackup)
	if err != nil {
		log.Printf("Failed to restore from backup on startup: %v", err)
		return
	}

	itemCount := js.PendingReviewData.Len()
	if itemCount == 0 {
		log.Printf("Restored backup %s, but it contains 0 items (empty backup)", latestBackup)
	} else {
		log.Printf("Successfully restored %d items from backup: %s", itemCount, latestBackup)
	}
}

func CheckServicesState(us *UserServices, js *JointServices) {
	validateConfiguration()
	logServiceInitialization(us, js)
	log.Println("Services state check completed")
}

func validateConfiguration() {
	if GetImageRoot() == "" {
		log.Println("WARNING: ImageRoot is not set")
	}
	if GetBackupDir() == "" {
		log.Println("WARNING: BackupDir is not set")
	}
}

func logServiceInitialization(us *UserServices, js *JointServices) {
	jobs := js.JobList.Jobs()
	if jobs == nil {
		log.Println("INFO: JobList initialized")
	}
	items := js.PendingReviewData.Items()
	if items == nil {
		log.Println("INFO: PendingReviewData initialized")
	}
	if us.CacheManager == nil {
		log.Println("INFO: CacheManager initialized")
	}
	if us.CurrentPageData.Len() == 0 {
		log.Println("INFO: CurrentPageData initialized")
	}
}
