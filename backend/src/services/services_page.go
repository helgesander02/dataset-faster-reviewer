package services

import (
	"backend/src/models_verify_viewer"
	"backend/src/utils"
	"log"
)

func (us *UserServices) SetCurrentPageData(jobName string, pageSize int) {
	log.Printf("[SetCurrentPageData] START - job: %s, pageSize: %d", jobName, pageSize)
	log.Printf("[SetCurrentPageData] BEFORE - CurrentPageData.JobName: %s, PageItems count: %d", us.CurrentPageData.JobName(), us.CurrentPageData.Len())

	root := GetImageRoot()
	jobData, _ := utils.ConcurrentJobDetailsScanner(root, jobName)
	us.FillJobNameToCurrentPageData(jobData.Name)

	log.Printf("[SetCurrentPageData] Scanned job: %s with %d datasets", jobData.Name, len(jobData.Datasets))

	for _, dataset := range jobData.Datasets {
		imageCount := dataset.GetImageLength()
		for i := 0; i < imageCount; i += pageSize {
			end := i + pageSize
			if end > imageCount {
				end = imageCount
			}
			us.AddPageToCurrentPageData(dataset.Name, dataset.Image[i:end])
		}
	}

	log.Printf("[SetCurrentPageData] AFTER - CurrentPageData.JobName: %s, PageItems count: %d", us.CurrentPageData.JobName(), us.CurrentPageData.Len())
}

func (us *UserServices) AddPageToCurrentPageData(datasetName string, datasetImage []models_verify_viewer.Image) {
	us.CurrentPageData.AddPage(datasetName, datasetImage)
}

func (us *UserServices) FillJobNameToCurrentPageData(jobName string) {
	us.CurrentPageData.SetJobName(jobName)
}

func (us *UserServices) GetCurrentPageData(jobName string) (*models_verify_viewer.Pages, bool) {
	if !us.currentPageDataExists(jobName) {
		log.Println("CurrentPageData does not exist for job:", jobName)
		return nil, false
	}
	return us.CurrentPageData, true
}

func (us *UserServices) GetImageCacheByPage(jobName string, pageIndex int) ([]string, []string) {
	return us.CurrentPageData.ImageNamesAt(pageIndex), us.CurrentPageData.ImagePathsAt(pageIndex)
}

func (us *UserServices) ClearCurrentPageData() {
	us.CurrentPageData.Clear()
}

func (us *UserServices) CurrentPageDataExists(jobName string) bool {
	return us.currentPageDataExists(jobName)
}

func (us *UserServices) currentPageDataExists(jobName string) bool {
	if us.CurrentPageData.JobName() == jobName && us.CurrentPageData.Len() > 0 {
		return true
	}
	return false
}

func (us *UserServices) RemoveImagesFromPageData(imagePaths []string) int {
	if len(imagePaths) == 0 {
		return 0
	}

	removedCount := us.CurrentPageData.RemoveImages(imagePaths)
	return removedCount
}
