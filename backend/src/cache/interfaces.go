package cache

import "backend/src/models"

type DataProvider interface {
	GetDatasets(jobName string) []string
	GetImages(jobName, datasetName string) []models.Image
}

type ImageProcessor interface {
	CompressToBase64(imagePath string) (string, error)
}
