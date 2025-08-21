package models_verify_viewer

import (
	"log"
	"os"
	"time"
)

type PendingReview struct {
	backupDir string
	Items     []PendingReviewItem `json:"items"`
}

type PendingReviewItem struct {
	JobName     string `json:"item_job_name"`
	DatasetName string `json:"item_dataset_name"`
	ImageName   string `json:"item_image_name"`
	ImagePath   string `json:"item_image_path"`
}

func NewPendingReview(backupDir string) PendingReview {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Printf("Failed to create backup directory: %v", err)
	}

	return PendingReview{
		backupDir: backupDir,
		Items:     NewPendingReviewItemSet(),
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

type BackupInfo struct {
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	ItemCount int       `json:"item_count"`
}
