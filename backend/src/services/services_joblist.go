package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
	"os"
)

func (dm *DataManager) GetParentDataAllJobs() []string {
	return dm.JobList.Jobs
}

func (dm *DataManager) GetParentDataAllDatasets(jobName string) []string {
	job, exists := dm.ConcurrentJobDetailsScanner(jobName)
	if !exists {
		log.Println("Failed to get job details")
		return []string{}
	}

	datasetNames := make([]string, len(job.Datasets))
	for i, ds := range job.Datasets {
		datasetNames[i] = ds.Name
	}

	return datasetNames
}

func (dm *DataManager) GetParentDataAllImages(jobName, datasetName string) []models_verify_viewer.Image {
	job, exists := dm.ConcurrentJobDetailsScanner(jobName)
	if !exists {
		log.Println("Failed to get job details")
		return []models_verify_viewer.Image{}
	}

	var images []models_verify_viewer.Image
	for _, ds := range job.Datasets {
		if ds.Name == datasetName {
			images = ds.Image
			break
		}
	}
	if images == nil {
		log.Println("Dataset not found")
		return []models_verify_viewer.Image{}
	}

	return images
}

func (dm *DataManager) ParentDataJobExists(jobName string) bool {
	for _, j := range dm.JobList.Jobs {
		if j == jobName {
			return true
		}
	}
	return false
}

func (dm *DataManager) ConcurrentJobScanner() {
	go utils.WatchJobs(dm.ImageRoot, &dm.JobList)
	log.Printf("Watchers initialized for root directory: %s", dm.ImageRoot)
}

func (dm *DataManager) ConcurrentJobDetailsScanner(jobName string) (models_verify_viewer.Job, bool) {
	jobPath := dm.ImageRoot + "/" + jobName
	_, err := os.Stat(jobPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Job directory does not exist: %s", jobPath)
		} else {
			log.Printf("Error accessing job directory: %v", err)
		}
		return models_verify_viewer.NewJob(jobName), false
	}

	job, exists := dm.GetJobCache(jobName)
	if exists {
		log.Printf("Job found in cache: %s", jobName)
		return job, true
	}

	job = utils.WatchJobDetails(dm.ImageRoot, jobName)
	log.Printf("Job not found in cache, scanning: %s", jobName)
	dm.MergeJobCache(job)
	return job, true
}

func (dm *DataManager) GetDatasets(jobName string) []string {
	return dm.GetParentDataAllDatasets(jobName)
}

func (dm *DataManager) GetImages(jobName, datasetName string) []models_verify_viewer.Image {
	return dm.GetParentDataAllImages(jobName, datasetName)
}
