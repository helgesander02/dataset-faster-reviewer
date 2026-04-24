package utils

import (
	"backend/src/models_verify_viewer"
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	watcherContext context.Context
	watcherCancel  context.CancelFunc
)

func ConcurrentJobScanner(ctx context.Context, root string, jobList *models_verify_viewer.JobList) {
	watcherContext, watcherCancel = context.WithCancel(ctx)
	go watchJobs(watcherContext, root, jobList)
	log.Printf("Job watcher initialized for root directory: %s", root)
}

func StopJobWatcher() {
	if watcherCancel != nil {
		watcherCancel()
		log.Println("Job watcher stopped gracefully")
	}
}

func watchJobs(ctx context.Context, root string, jobList *models_verify_viewer.JobList) {
	watcher, err := createWatcher(root)
	if err != nil {
		return
	}
	defer watcher.Close()

	performInitialScan(root, jobList)
	log.Printf("Watching directory: %s", root)

	monitorFileSystemEvents(ctx, watcher, root, jobList)
}

func createWatcher(root string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create watcher: %v", err)
		return nil, err
	}

	if err := watcher.Add(root); err != nil {
		log.Printf("Failed to add watcher for %s: %v", root, err)
		watcher.Close()
		return nil, err
	}

	return watcher, nil
}

func performInitialScan(root string, jobList *models_verify_viewer.JobList) {
	jobs := scanJobs(root)
	jobList.Replace(jobs)
	log.Printf("Initial scan found %d jobs", len(jobs))
}

func monitorFileSystemEvents(ctx context.Context, watcher *fsnotify.Watcher, root string, jobList *models_verify_viewer.JobList) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Job watcher context cancelled, shutting down...")
			return

		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Watcher events channel closed")
				return
			}
			handleFileSystemEvent(event, root, jobList)

		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("Watcher errors channel closed")
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func handleFileSystemEvent(event fsnotify.Event, root string, jobList *models_verify_viewer.JobList) {
	if !isJobLevelChange(root, event.Name) {
		return
	}

	if !isRelevantOperation(event.Op) {
		return
	}

	log.Printf("Detected job level change: %s [%s]", event.Name, event.Op)
	refreshJobList(root, jobList)
}

func isRelevantOperation(op fsnotify.Op) bool {
	return op&fsnotify.Write == fsnotify.Write ||
		op&fsnotify.Create == fsnotify.Create ||
		op&fsnotify.Remove == fsnotify.Remove
}

func refreshJobList(root string, jobList *models_verify_viewer.JobList) {
	jobs := scanJobs(root)
	jobList.Replace(jobs)
	log.Printf("Refreshed job list: %d jobs", len(jobs))
}

func scanJobs(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("Error reading root directory: %v", err)
		return []string{}
	}

	return extractJobNames(entries)
}

func extractJobNames(entries []os.DirEntry) []string {
	var jobs []string
	for _, entry := range entries {
		if entry.IsDir() {
			jobs = append(jobs, entry.Name())
		}
	}
	return jobs
}

func isJobLevelChange(root string, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}

	return !filepath.IsAbs(rel) && (filepath.Dir(rel) == "." || rel == ".")
}
