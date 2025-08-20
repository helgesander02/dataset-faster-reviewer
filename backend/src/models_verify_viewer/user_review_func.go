package models_verify_viewer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func (pr_old *PendingReview) MergePendingReviewItems(pr_new PendingReview) {
	log.Printf("Merging %d new items into existing %d items", len(pr_new.Items), len(pr_old.Items))

	newItemsMap := make(map[string]bool)
	for _, newitem := range pr_new.Items {
		key := newitem.JobName + "|" + newitem.DatasetName + "|" + newitem.ImageName
		newItemsMap[key] = true
	}

	filteredOldItems := NewPendingReviewItemSet()
	for _, oldItem := range pr_old.Items {
		key := oldItem.JobName + "|" + oldItem.DatasetName + "|" + oldItem.ImageName
		if newItemsMap[key] {
			filteredOldItems = append(filteredOldItems, oldItem)
		}
	}
	pr_old.Items = filteredOldItems

	for _, newitem := range pr_new.Items {
		found := false
		for _, oldItem := range pr_old.Items {
			if newitem.JobName == oldItem.JobName && newitem.DatasetName == oldItem.DatasetName && newitem.ImageName == oldItem.ImageName {
				found = true
				break
			}
		}
		if !found {
			pr_old.Items = append(pr_old.Items, newitem)
		}
	}
}

func (pr *PendingReview) ClearPendingReviewItems() {
	pr.Items = []PendingReviewItem{}
}

// Backup
const (
	MaxBackupCount = 10
)

func (bm *BackupManager) CreateBackup(pendingReview PendingReview) error {
	timestamp := time.Now()
	filename := fmt.Sprintf("pending_review_%s.json", timestamp.Format("20060102_150405"))
	backupPath := filepath.Join(bm.backupDir, filename)

	jsonData, err := json.MarshalIndent(pendingReview, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pending review data: %v", err)
	}

	if err := os.WriteFile(backupPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	log.Printf("Created backup: %s with %d items", filename, len(pendingReview.Items))

	if err := bm.cleanupOldBackups(); err != nil {
		log.Printf("Warning: failed to cleanup old backups: %v", err)
	}

	return nil
}

func (bm *BackupManager) cleanupOldBackups() error {
	files, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backupFiles []os.DirEntry
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "pending_review_") && strings.HasSuffix(file.Name(), ".json") {
			backupFiles = append(backupFiles, file)
		}
	}

	sort.Slice(backupFiles, func(i, j int) bool {
		infoI, _ := backupFiles[i].Info()
		infoJ, _ := backupFiles[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	if len(backupFiles) > MaxBackupCount {
		for i := MaxBackupCount; i < len(backupFiles); i++ {
			oldBackupPath := filepath.Join(bm.backupDir, backupFiles[i].Name())
			if err := os.Remove(oldBackupPath); err != nil {
				log.Printf("Failed to remove old backup %s: %v", backupFiles[i].Name(), err)
			} else {
				log.Printf("Removed old backup: %s", backupFiles[i].Name())
			}
		}
	}

	return nil
}

func (bm *BackupManager) ListBackups() ([]BackupInfo, error) {
	files, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backups []BackupInfo
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "pending_review_") && strings.HasSuffix(file.Name(), ".json") {
			itemCount := 0
			backupPath := filepath.Join(bm.backupDir, file.Name())
			if data, err := os.ReadFile(backupPath); err == nil {
				var pr PendingReview
				if err := json.Unmarshal(data, &pr); err == nil {
					itemCount = len(pr.Items)
				}
			}

			info, _ := file.Info()
			backups = append(backups, BackupInfo{
				Filename:  file.Name(),
				Timestamp: info.ModTime(),
				ItemCount: itemCount,
			})
		}
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	return backups, nil
}

func (bm *BackupManager) RestoreFromBackup(filename string) (PendingReview, error) {
	backupPath := filepath.Join(bm.backupDir, filename)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return PendingReview{}, fmt.Errorf("backup file not found: %s", filename)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return PendingReview{}, fmt.Errorf("failed to read backup file: %v", err)
	}

	var pendingReview PendingReview
	if err := json.Unmarshal(data, &pendingReview); err != nil {
		return PendingReview{}, fmt.Errorf("failed to unmarshal backup data: %v", err)
	}

	log.Printf("Restored from backup: %s with %d items", filename, len(pendingReview.Items))
	return pendingReview, nil
}

func (bm *BackupManager) GetBackupPath() string {
	return bm.backupDir
}
