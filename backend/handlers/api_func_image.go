// handlers/api_func_image.go
package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// @Summary      Get folder structure
// @Description  Returns the folder structure data
// @Tags         folder
// @Produce      json
// @Success      200  {object}  interface{}
// @Router       /api/folder-structure [get]
func (handle *Handle) FolderStructureHandler(c *gin.Context) {
	data := handle.DM.ParentData
	c.JSON(http.StatusOK, data)
}

// @Summary      Get all jobs
// @Description  Returns a list of all job names
// @Tags         jobs
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/getJobs [get]
func (handle *Handle) GetJobs(c *gin.Context) {
	jobNames := handle.DM.GetParentDataAllJobs()
	c.JSON(http.StatusOK, gin.H{
		"total_jobs": len(jobNames),
		"job_names":  jobNames,
	})
}

// @Summary      Get all datasets for a job
// @Description  Returns a list of all dataset names for a given job
// @Tags         datasets
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/getDatasets [get]
func (handle *Handle) GetDatasets(c *gin.Context) {
	jobName := c.Query("job")
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job parameter"})
		return
	}
	
	if !handle.DM.ParentDataJobExists(jobName) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	
	datasetNames := handle.DM.GetParentDataAllDatasets(jobName)
	c.JSON(http.StatusOK, gin.H{
		"total_datasets": len(datasetNames),
		"dataset_names":  datasetNames,
	})
}

// @Summary      Get all images for a dataset
// @Description  Returns a list of all images for a given job and dataset
// @Tags         images
// @Produce      json
// @Param        job      query  string  true  "Job name"
// @Param        dataset  query  string  true  "Dataset name"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/getImages [get]
func (handle *Handle) GetImages(c *gin.Context) {
	jobName := c.Query("job")
	datasetName := c.Query("dataset")
	
	if jobName == "" || datasetName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job or dataset parameter"})
		return
	}
	
	if !handle.DM.ParentDataJobExists(jobName) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	
	images := handle.DM.GetParentDataAllImages(jobName, datasetName)
	c.JSON(http.StatusOK, gin.H{
		"total_images": len(images),
		"images":       images,
	})
}

// @Summary      Get base64 images for a dataset page
// @Description  Returns base64-encoded images for a given job, dataset, pageIndex, and pageNumber
// @Tags         images
// @Produce      json
// @Param        job        query  string  true  "Job name"
// @Param        dataset    query  string  true  "Dataset name (not used in cache, but kept for compatibility)"
// @Param        pageIndex  query  int     true  "Page index"
// @Param        pageNumber query  int     true  "Page size (for initialization only)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getBase64Images [get]
func (handle *Handle) GetBase64Images(c *gin.Context) {
	jobName := c.Query("job")
	//datasetName := c.Query("dataset") // 保留但不使用，為了向後兼容
	pageIndexStr := c.Query("pageIndex")
	pageNumberStr := c.Query("pageNumber")
	
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job parameter"})
		return
	}
	
	if !handle.DM.ParentDataJobExists(jobName) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	
	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageIndex parameter"})
		return
	}
	
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageNumber parameter"})
		return
	}
	
	// 檢查快取是否存在，如果不存在則初始化
	if !handle.DM.ImageCacheJobExists(jobName) {
		if err := handle.DM.InitializeImagesCache(jobName, pageNumber); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to initialize cache: " + err.Error(),
			})
			return
		}
	}
	
	// 使用新的 cache API (只需要 jobName 和 pageIndex)
	cachedImages, err := handle.DM.GetImagesCache(jobName, pageIndex)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to get images: " + err.Error(),
		})
		return
	}
	
	// 取得最大頁數
	maxPage, err := handle.DM.GetJobMaxPages(jobName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get max pages: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"max_page": maxPage,
		"images":   cachedImages,
	})
}

// @Summary      Get all pages for a job
// @Description  Returns all page details for a given job
// @Tags         pages
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getAllPages [get]
func (handle *Handle) GetAllPages(c *gin.Context) {
	jobName := c.Query("job")
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job parameter"})
		return
	}
	
	if !handle.DM.ImageCacheJobExists(jobName) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job cache not found"})
		return
	}
	
	pages, err := handle.DM.GetJobPageDetail(jobName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get page details: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total_pages": len(pages),
		"pages":       pages,
	})
}
