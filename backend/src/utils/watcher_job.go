package utils

import (
	"backend/src/models_verify_viewer"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var jobMu sync.Mutex

func WatchJobs(root string, JobList *models_verify_viewer.JobList) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	jobMu.Lock()
	*JobList = *scanJobs(root)
	jobMu.Unlock()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if isJobLevelChange(root, event.Name) &&
					(event.Op&fsnotify.Write == fsnotify.Write ||
						event.Op&fsnotify.Create == fsnotify.Create ||
						event.Op&fsnotify.Remove == fsnotify.Remove) {
					log.Println("Detected job level change:", event)

					jobMu.Lock()
					*JobList = *scanJobs(root)
					jobMu.Unlock()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(root)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func scanJobs(root string) *models_verify_viewer.JobList {
	parent := models_verify_viewer.NewJobList()
	jobs, err := os.ReadDir(root)
	if err != nil {
		log.Printf("Error reading root directory: %v", err)
		return &parent
	}

	for _, job := range jobs {
		if !job.IsDir() {
			continue
		}
		parent.Jobs = append(parent.Jobs, job.Name())
	}

	return &parent
}

func isJobLevelChange(root string, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && filepath.Dir(rel) == "." || rel == "."
}
