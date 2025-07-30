package services

import "backend/src/models"

// DataManager 實現 DataProvider 介面
func (dm *DataManager) GetDatasets(jobName string) []string {
	return dm.GetParentDataAllDatasets(jobName)
}

func (dm *DataManager) GetImages(jobName, datasetName string) []models.Image {
	return dm.GetParentDataAllImages(jobName, datasetName)
}
