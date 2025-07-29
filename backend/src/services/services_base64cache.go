package services

import (
	"backend/src/models"
	"log"
)

// 移除舊的 GetAllPageDetail 函數，因為它需要 jobName 參數
// 使用 GetJobPageDetail 替代

func (dm *DataManager) InitialImagesCache(jobName string, pageNumber int) {
	log.Printf("Initializing images cache for job: %s with page size: %d", jobName, pageNumber)
	
	getDatasetsFn := func(job string) []string {
		return dm.GetParentDataAllDatasets(job)
	}
	
	getImagesFn := func(job, dataset string) []models.Image {
		return dm.GetParentDataAllImages(job, dataset)
	}
	
	dm.CacheManager.InitialImagesCache(jobName, pageNumber, getDatasetsFn, getImagesFn)
}

func (dm *DataManager) GetImagesCache(jobName string, pageIndex int, pageNumber int) (models.ImageItems, bool) {
	getDatasetsFn := func(job string) []string {
		return dm.GetParentDataAllDatasets(job)
	}
	
	getImagesFn := func(job, dataset string) []models.Image {
		return dm.GetParentDataAllImages(job, dataset)
	}
	
	return dm.CacheManager.GetImagesCache(jobName, pageIndex, pageNumber, getDatasetsFn, getImagesFn)
}

func (dm *DataManager) ClearImagesCache(jobName string) bool {
	dm.CacheManager.ClearImagesCache(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
	return true
}

func (dm *DataManager) ImageCacheJobExists(jobName string) bool {
	return dm.CacheManager.ImageCacheJobExists(jobName)
}

// 新增：取得特定 job 的頁面詳細資訊
func (dm *DataManager) GetJobPageDetail(jobName string) map[int]string {
	return dm.CacheManager.GetAllPageDetail(jobName)
}

// 新增：取得指定 job 的最大頁數
func (dm *DataManager) GetJobMaxPages(jobName string) int {
	return dm.CacheManager.GetJobMaxPages(jobName)
}
