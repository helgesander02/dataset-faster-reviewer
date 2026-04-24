package utils

import (
	"backend/src/models_verify_viewer"
	"log"
	"os"
	"path/filepath"
)

const (
	imageSubdirectory = "image"
	labelSubdirectory = "label"
	jpgExtension      = ".jpg"
	jsonExtension     = ".json"
)

func ConcurrentJobDetailsScanner(root string, jobName string) (models_verify_viewer.Job, bool) {
	jobPath := filepath.Join(root, jobName)

	if !jobDirectoryExists(jobPath) {
		return createEmptyJob(jobName), false
	}

	jobData := watchJobDetails(root, jobName)
	return jobData, true
}

func jobDirectoryExists(jobPath string) bool {
	_, err := os.Stat(jobPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Job directory does not exist: %s", jobPath)
		} else {
			log.Printf("Error accessing job directory: %v", err)
		}
		return false
	}
	return true
}

func createEmptyJob(jobName string) models_verify_viewer.Job {
	job := models_verify_viewer.NewJob()
	job.FillJobName(jobName)
	return job
}

func watchJobDetails(root string, jobName string) models_verify_viewer.Job {
	jobData := models_verify_viewer.NewJob()
	jobData.FillJobName(jobName)

	jobPath := filepath.Join(root, jobName)
	if !jobDirectoryExists(jobPath) {
		return jobData
	}

	datasets := readDatasets(jobPath, jobName)
	jobData.Datasets = processDatasets(jobPath, datasets)

	return jobData
}

func readDatasets(jobPath, jobName string) []os.DirEntry {
	datasets, err := os.ReadDir(jobPath)
	if err != nil {
		log.Printf("Error reading datasets for job %s: %v", jobName, err)
		return []os.DirEntry{}
	}
	return datasets
}

func processDatasets(jobPath string, datasets []os.DirEntry) []models_verify_viewer.Dataset {
	var processedDatasets []models_verify_viewer.Dataset
	for _, dataset := range datasets {
		if dataset.IsDir() {
			datasetData := scanDataset(jobPath, dataset)
			processedDatasets = append(processedDatasets, datasetData)
		}
	}
	return processedDatasets
}

func scanDataset(jobPath string, dataset os.DirEntry) models_verify_viewer.Dataset {
	datasetData := models_verify_viewer.NewDataset()
	datasetData.FillDatasetName(dataset.Name())

	datasetPath := filepath.Join(jobPath, dataset.Name())
	metaDirectories := readMetaDirectories(datasetPath, dataset.Name())

	processMetaDirectories(datasetPath, metaDirectories, &datasetData)

	return datasetData
}

func readMetaDirectories(datasetPath, datasetName string) []os.DirEntry {
	metas, err := os.ReadDir(datasetPath)
	if err != nil {
		log.Printf("Error reading meta directories for dataset %s: %v", datasetName, err)
		return []os.DirEntry{}
	}
	return metas
}

func processMetaDirectories(datasetPath string, metas []os.DirEntry, datasetData *models_verify_viewer.Dataset) {
	for _, meta := range metas {
		if meta.IsDir() {
			metaPath := filepath.Join(datasetPath, meta.Name())
			scanMeta(metaPath, meta.Name(), datasetData)
		}
	}
}

func scanMeta(metaPath, metaName string, datasetData *models_verify_viewer.Dataset) {
	switch metaName {
	case imageSubdirectory:
		scanImagesDirectory(metaPath, datasetData)
	case labelSubdirectory:
		scanLabelsDirectory(metaPath, datasetData)
	}
}

func scanImagesDirectory(metaPath string, datasetData *models_verify_viewer.Dataset) {
	images, err := os.ReadDir(metaPath)
	if err != nil {
		log.Printf("Error reading images: %v", err)
		return
	}
	scanImages(images, metaPath, datasetData)
}

func scanLabelsDirectory(metaPath string, datasetData *models_verify_viewer.Dataset) {
	labels, err := os.ReadDir(metaPath)
	if err != nil {
		log.Printf("Error reading labels: %v", err)
		return
	}
	scanLabels(labels, metaPath, datasetData)
}

func scanImages(images []os.DirEntry, metaPath string, datasetData *models_verify_viewer.Dataset) {
	capacity := countValidFiles(images, jpgExtension)
	datasetData.Image = make([]models_verify_viewer.Image, 0, capacity)

	for _, image := range images {
		if isValidImageFile(image) {
			imagePath := filepath.Join(metaPath, image.Name())
			datasetData.Image = append(datasetData.Image, models_verify_viewer.NewImage(image.Name(), imagePath))
		}
	}
}

func scanLabels(labels []os.DirEntry, metaPath string, datasetData *models_verify_viewer.Dataset) {
	capacity := countValidFiles(labels, jsonExtension)
	datasetData.Label = make([]models_verify_viewer.Label, 0, capacity)

	for _, label := range labels {
		if isValidLabelFile(label) {
			labelPath := filepath.Join(metaPath, label.Name())
			datasetData.Label = append(datasetData.Label, models_verify_viewer.NewLabel(label.Name(), labelPath))
		}
	}
}

func countValidFiles(files []os.DirEntry, extension string) int {
	count := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == extension {
			count++
		}
	}
	return count
}

func isValidImageFile(file os.DirEntry) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == jpgExtension
}

func isValidLabelFile(file os.DirEntry) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == jsonExtension
}
