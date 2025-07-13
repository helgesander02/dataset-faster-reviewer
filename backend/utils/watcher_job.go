package utils

import (
	"log"
	"os"
	"sync"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
	"backend/models"
)

var jobMu sync.Mutex

func WatchJobs(root string, ParentData *models.Parent) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    jobMu.Lock()
    *ParentData = *scanJobs(root)
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
                    *ParentData = *scanJobs(root)
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


func scanJobs(root string) *models.Parent {
    parent := models.NewParentData()
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
