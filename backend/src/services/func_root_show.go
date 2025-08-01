package services

import (
	"backend/src/models_verify_viewer"
	"log"
)

func (us *UserServices) GetDatasets(jobName string) []string {
	return us.GetAllDatasets(jobName)
}

func (us *UserServices) GetImages(jobName, datasetName string) []models_verify_viewer.Image {
	return us.GetAllImages(jobName, datasetName)
}

func (us *UserServices) GetAllDatasets(jobName string) []string {
	job, exists := us.ConcurrentJobDetailsScanner(jobName)
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

func (us *UserServices) GetAllImages(jobName, datasetName string) []models_verify_viewer.Image {
	job, exists := us.ConcurrentJobDetailsScanner(jobName)
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
