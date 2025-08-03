package utils

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var jobMu sync.Mutex

func ConcurrentJobScanner(root string, job_list *[]string) {
	go watchJobs(root, job_list)
	log.Printf("Watchers initialized for root directory: %s", root)
}

func watchJobs(root string, job_list *[]string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 初始化掃描
	jobMu.Lock()
	jobs := scanJobs(root)
	*job_list = jobs
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
					jobs := scanJobs(root)
					*job_list = jobs
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

func scanJobs(root string) []string {
	var jobs []string
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("Error reading root directory: %v", err)
		return jobs
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		jobs = append(jobs, entry.Name())
	}

	return jobs
}

func isJobLevelChange(root string, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && filepath.Dir(rel) == "." || rel == "."
}
