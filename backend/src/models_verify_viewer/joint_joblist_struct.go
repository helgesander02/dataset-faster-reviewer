package models_verify_viewer

import "sync"

type JobList struct {
	jobs []string
	mu   sync.RWMutex
}

func NewJobList() *JobList {
	return &JobList{
		jobs: make([]string, 0),
	}
}
