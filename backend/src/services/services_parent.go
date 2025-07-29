package services

import (
	"backend/src/models"
	"log"
)

func (dm *DataManager) GetParentDataAllJobs() []string {
	return dm.ParentData.Jobs
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

func (dm *DataManager) GetParentDataAllImages(jobName, datasetName string) []models.Image {
	job, exists := dm.ConcurrentJobDetailsScanner(jobName)
	if !exists {
		log.Println("Failed to get job details")
		return []models.Image{}
	}

	var images []models.Image
	for _, ds := range job.Datasets {
		if ds.Name == datasetName {
			images = ds.Image
			break
		}
	}
	if images == nil {
		log.Println("Dataset not found")
		return []models.Image{}
	}

	return images
}

func (dm *DataManager) ParentDataJobExists(jobName string) bool {
	for _, j := range dm.ParentData.Jobs {
		if j == jobName {
			return true
		}
	}
	return false
}
