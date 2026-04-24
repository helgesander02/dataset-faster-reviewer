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

const (
	MaxBackupCount        = 10
	backupFilenameFormat  = "pending_review_%s.json"
	backupTimestampFormat = "20060102_150405"
	backupFilePermissions = 0644
)

func (pr *PendingReview) GetLatestBackup(backupDir string) (string, error) {
	backups, err := pr.ListBackups(backupDir)
	if err != nil {
		return "", fmt.Errorf("failed to list backups: %v", err)
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found")
	}

	return backups[0].Filename, nil
}

func (pr *PendingReview) RestoreFromBackup(backupDir string, filename string) error {
	backupPath := filepath.Join(backupDir, filename)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", filename)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %v", err)
	}

	var temp struct {
		Items []PendingReviewItem `json:"items"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal backup data: %v", err)
	}

	pr.mu.Lock()
	pr.items = temp.Items
	pr.mu.Unlock()

	return nil
}

func (pr *PendingReview) CreateBackup(backupDir string) error {
	backupPath, err := pr.prepareBackupFile(backupDir)
	if err != nil {
		return err
	}

	if err := pr.writeBackupFile(backupPath); err != nil {
		return err
	}

	pr.cleanupOldBackups(backupDir)
	return nil
}

func (pr *PendingReview) prepareBackupFile(backupDir string) (string, error) {
	ensureBackupDirectoryExists(backupDir)
	timestamp := time.Now()
	filename := fmt.Sprintf(backupFilenameFormat, timestamp.Format(backupTimestampFormat))
	return filepath.Join(backupDir, filename), nil
}

func (pr *PendingReview) writeBackupFile(backupPath string) error {
	pr.mu.RLock()

	temp := struct {
		Items []PendingReviewItem `json:"items"`
	}{
		Items: pr.items,
	}
	jsonData, err := json.MarshalIndent(temp, "", "  ")
	pr.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal pending review data: %v", err)
	}

	if err := os.WriteFile(backupPath, jsonData, backupFilePermissions); err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	return nil
}

func (pr *PendingReview) cleanupOldBackups(backupDir string) {
	backupFiles, err := pr.getBackupFiles(backupDir)
	if err != nil {
		log.Printf("Warning: failed to cleanup old backups: %v", err)
		return
	}

	pr.sortBackupFilesByDate(backupFiles)
	pr.removeExcessBackups(backupDir, backupFiles)
}

func (pr *PendingReview) getBackupFiles(backupDir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backupFiles []os.DirEntry
	for _, file := range files {
		if isBackupFile(file.Name()) {
			backupFiles = append(backupFiles, file)
		}
	}

	return backupFiles, nil
}

func isBackupFile(filename string) bool {
	return strings.HasPrefix(filename, "pending_review_") && strings.HasSuffix(filename, ".json")
}

func (pr *PendingReview) sortBackupFilesByDate(backupFiles []os.DirEntry) {
	sort.Slice(backupFiles, func(i, j int) bool {
		infoI, _ := backupFiles[i].Info()
		infoJ, _ := backupFiles[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})
}

func (pr *PendingReview) removeExcessBackups(backupDir string, backupFiles []os.DirEntry) {
	if len(backupFiles) <= MaxBackupCount {
		return
	}

	for i := MaxBackupCount; i < len(backupFiles); i++ {
		pr.removeBackupFile(backupDir, backupFiles[i])
	}
}

func (pr *PendingReview) removeBackupFile(backupDir string, file os.DirEntry) {
	backupPath := filepath.Join(backupDir, file.Name())
	if err := os.Remove(backupPath); err != nil {
		log.Printf("Failed to remove old backup %s: %v", file.Name(), err)
	}
}

func (pr *PendingReview) ListBackups(backupDir string) ([]BackupInfo, error) {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	backups := pr.collectBackupInfo(backupDir, files)
	pr.sortBackupsByTimestamp(backups)

	return backups, nil
}

func (pr *PendingReview) collectBackupInfo(backupDir string, files []os.DirEntry) []BackupInfo {
	var backups []BackupInfo
	for _, file := range files {
		if !isBackupFile(file.Name()) {
			continue
		}

		backupInfo := pr.createBackupInfo(backupDir, file)
		backups = append(backups, backupInfo)
	}
	return backups
}

func (pr *PendingReview) createBackupInfo(backupDir string, file os.DirEntry) BackupInfo {
	info, _ := file.Info()
	itemCount := pr.getBackupItemCount(backupDir, file.Name())

	return BackupInfo{
		Filename:  file.Name(),
		Timestamp: info.ModTime(),
		ItemCount: itemCount,
	}
}

func (pr *PendingReview) getBackupItemCount(backupDir string, filename string) int {
	backupPath := filepath.Join(backupDir, filename)
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return 0
	}

	var temp struct {
		Items []PendingReviewItem `json:"items"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return 0
	}

	return len(temp.Items)
}

func (pr *PendingReview) sortBackupsByTimestamp(backups []BackupInfo) {
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})
}
