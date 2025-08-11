package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
)

func (us *UserServices) SetCurrentPageData(jobName string, pageSize int) {
	job, _ := utils.ConcurrentJobDetailsScanner(ImageRoot, jobName)
	us.CurrentPageData.FillJobName(job.Name)

	for _, dataset := range job.Datasets {
		imageCount := dataset.GetImageLength()
		for i := 0; i < imageCount; i += pageSize {
			end := i + pageSize
			if end > imageCount {
				end = imageCount
			}
			us.CurrentPageData.AppendPageItems(dataset.Name, dataset.Image[i:end])
		}
	}
}

func (us *UserServices) GetCurrentPageData(jobName string) (models_verify_viewer.Pages, bool) {
	if !us.CurrentPageDataExists(jobName) {
		log.Println("CurrentPageData does not exist for job:", jobName)
		return models_verify_viewer.NewPages(), false
	}
	return us.CurrentPageData, true
}

func (us UserServices) GetImageCacheByPage(jobName string, pageIndex int) ([]string, []string) {
	return us.CurrentPageData.GetPageItemAllImageNameByIndex(pageIndex), us.CurrentPageData.GetPageItemAllImagePathByIndex(pageIndex)
}

func (us *UserServices) ClearCurrentPageData() {
	us.CurrentPageData.ClearPages()
}

func (us *UserServices) CurrentPageDataExists(jobName string) bool {
	if us.CurrentPageData.JobName == jobName && us.CurrentPageData.GetPageItemsLength() > 0 {
		return true
	}
	return false
}
