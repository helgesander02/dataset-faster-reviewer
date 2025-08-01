package services

import (
	"backend/src/models_verify_viewer"
	"log"
)

func (us *UserServices) InitializeImagesCache(jobName string, pageSize int) error {
	log.Printf("Initializing images cache for job: %s with page size: %d", jobName, pageSize)
	return us.CacheManager.InitializeImagesCache(jobName, pageSize)
}

func (us *UserServices) GetImagesCache(jobName string, pageIndex int) (models_verify_viewer.ImageItems, error) {
	return us.CacheManager.GetImagesCache(jobName, pageIndex)
}

func (us *UserServices) ClearImagesCache(jobName string) error {
	us.CacheManager.ClearImagesCache(jobName)
	log.Printf("Cleared images cache for job: %s", jobName)
	return nil
}

func (us *UserServices) ImageCacheJobExists(jobName string) bool {
	return us.CacheManager.ImageCacheJobExists(jobName)
}

func (us *UserServices) GetJobPageDetail(jobName string) (map[int]string, error) {
	return us.CacheManager.GetAllPageDetail(jobName)
}

func (us *UserServices) GetJobMaxPages(jobName string) (int, error) {
	return us.CacheManager.GetJobMaxPages(jobName)
}

func (us *UserServices) GetJobCache(jobName string) (models_verify_viewer.Job, bool) {
	return us.CacheManager.GetJobCache(jobName)
}

func (us *UserServices) MergeJobCache(job models_verify_viewer.Job) {
	us.CacheManager.SetJobCache(job)
}
