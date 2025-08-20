package models_verify_viewer

import (
	"log"
	"os"
	"time"
)

type PendingReview struct {
	Items []PendingReviewItem `json:"items"`
}

type PendingReviewItem struct {
	JobName     string `json:"item_job_name"`
	DatasetName string `json:"item_dataset_name"`
	ImageName   string `json:"item_image_name"`
	ImagePath   string `json:"item_image_path"`
}

func NewPendingReview() PendingReview {
	return PendingReview{
		Items: NewPendingReviewItemSet(),
	}
}

func NewPendingReviewItem() PendingReviewItem {
	return PendingReviewItem{
		JobName:     "",
		DatasetName: "",
		ImageName:   "",
		ImagePath:   "",
	}
}

func NewPendingReviewItemSet() []PendingReviewItem {
	return []PendingReviewItem{}
}

func NewPendingReviewItemSetByLenght(lenght int) []PendingReviewItem {
	return make([]PendingReviewItem, 0, lenght)
}

type BackupManager struct {
	backupDir string
}

type BackupInfo struct {
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	ItemCount int       `json:"item_count"`
}

func NewBackupManager(backupDir string) *BackupManager {
	bm := &BackupManager{
		backupDir: backupDir,
	}
	if err := os.MkdirAll(bm.backupDir, 0755); err != nil {
		log.Printf("Failed to create backup directory: %v", err)
	}

	return bm
}
