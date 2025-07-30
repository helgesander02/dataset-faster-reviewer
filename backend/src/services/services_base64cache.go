package services

import (
	"backend/src/models"
	"log"
)

func (dm *DataManager) InitializeImagesCache(jobName string, pageSize int) error {
	log.Printf("Initializing images cache for job: %s with page size: %d", jobName, pageSize)
	return dm.CacheManager.InitializeImagesCache(jobName, pageSize)
}

func (dm *DataManager) GetImagesCache(jobName string, pageIndex int) (models.ImageItems, error) {
	return dm.CacheManager.GetImagesCache(jobName, pageIndex)
}

func (dm *DataManager) ClearImagesCache(jobName string) error {
	dm.CacheManager.ClearImagesCache(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
	return nil
}

func (dm *DataManager) ImageCacheJobExists(jobName string) bool {
	return dm.CacheManager.ImageCacheJobExists(jobName)
}

func (dm *DataManager) GetJobPageDetail(jobName string) (map[int]string, error) {
	return dm.CacheManager.GetAllPageDetail(jobName)
}

func (dm *DataManager) GetJobMaxPages(jobName string) (int, error) {
	return dm.CacheManager.GetJobMaxPages(jobName)
}
