package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
	"os"
)

func (js *JointServices) ConcurrentJobScanner() {
	go utils.WatchJobs(ImageRoot, &js.JobList)
	log.Printf("Watchers initialized for root directory: %s", ImageRoot)
}

func (us *UserServices) ConcurrentJobDetailsScanner(jobName string) (models_verify_viewer.Job, bool) {
	jobPath := ImageRoot + "/" + jobName
	_, err := os.Stat(jobPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Job directory does not exist: %s", jobPath)
		} else {
			log.Printf("Error accessing job directory: %v", err)
		}
		return models_verify_viewer.NewJob(jobName), false
	}

	job, exists := us.GetJobCache(jobName)
	if exists {
		log.Printf("Job found in cache: %s", jobName)
		return job, true
	}

	job = utils.WatchJobDetails(ImageRoot, jobName)
	log.Printf("Job not found in cache, scanning: %s", jobName)
	us.MergeJobCache(job)
	return job, true
}
