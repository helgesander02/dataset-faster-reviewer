package models_verify_viewer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type PendingReview struct {
	items []PendingReviewItem
	mu    sync.RWMutex
}

type PendingReviewItem struct {
	JobName     string `json:"item_job_name"`
	DatasetName string `json:"item_dataset_name"`
	ImageName   string `json:"item_image_name"`
	ImagePath   string `json:"item_image_path"`
}

func NewPendingReview() *PendingReview {
	return &PendingReview{
		items: make([]PendingReviewItem, 0),
	}
}

func NewPendingReviewItem(jobName, datasetName, imageName, imagePath string) (PendingReviewItem, error) {
	if jobName == "" {
		return PendingReviewItem{}, fmt.Errorf("jobName cannot be empty")
	}
	if datasetName == "" {
		return PendingReviewItem{}, fmt.Errorf("datasetName cannot be empty")
	}
	if imageName == "" {
		return PendingReviewItem{}, fmt.Errorf("imageName cannot be empty")
	}
	if imagePath == "" {
		return PendingReviewItem{}, fmt.Errorf("imagePath cannot be empty")
	}

	return PendingReviewItem{
		JobName:     jobName,
		DatasetName: datasetName,
		ImageName:   imageName,
		ImagePath:   imagePath,
	}, nil
}

func (item PendingReviewItem) Key() string {
	return fmt.Sprintf("%s|%s|%s", item.JobName, item.DatasetName, item.ImageName)
}

// MarshalJSON implements json.Marshaler interface for PendingReview
// This allows proper JSON serialization even though items is a private field
func (pr *PendingReview) MarshalJSON() ([]byte, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	// Create a temporary struct with public field for JSON serialization
	temp := struct {
		Items []PendingReviewItem `json:"items"`
	}{
		Items: pr.items,
	}

	return json.Marshal(temp)
}

func ensureBackupDirectoryExists(backupDir string) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Printf("Failed to create backup directory: %v", err)
	}
}
