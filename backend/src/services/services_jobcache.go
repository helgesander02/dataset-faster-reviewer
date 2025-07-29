package services

import (
	"backend/src/models"
)

func (dm *DataManager) GetJobCache(jobName string) (models.Job, bool) {
	return dm.CacheManager.GetJobCache(jobName)
}

func (dm *DataManager) MergeJobCache(job models.Job) {
	dm.CacheManager.SetJobCache(job)
}
