package models_verify_viewer

import "time"

type BackupInfo struct {
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	ItemCount int       `json:"item_count"`
}
