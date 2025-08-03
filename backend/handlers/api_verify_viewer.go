// handlers/api_func_image.go
package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary      Get all jobs
// @Description  Returns a list of all job names
// @Tags         jobs
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/getJobs [get]
func (handle *Handle) GetJobs(c *gin.Context) {
	jobNames := handle.JointServices.GetJobList()
	c.JSON(http.StatusOK, gin.H{
		"total_jobs": len(jobNames),
		"job_names":  jobNames,
	})
}

// @Summary      Set all pages for a job
// @Description  Sets all page details for a given job
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Param        pages  body  []string  true  "List of page details"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/setAllPages [post]
func (handle *Handle) SetAllPageDetails(c *gin.Context) {
	jobName := c.Query("job")
	pageSize := c.Query("pageSize")
	if jobName == "" || pageSize == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job or pageSize parameter"})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	if !handle.UserServices.CurrentPageDataExists(jobName) {
		handle.UserServices.CacheManager.ClearImageCacheStore(jobName)
		handle.UserServices.ClearCurrentPageData()
		handle.UserServices.SetCurrentPageData(jobName, pageSizeInt)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "All pages set successfully for job: " + jobName,
		"job":       jobName,
		"page_size": pageSizeInt,
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
func (handle *Handle) GetAllPageDetails(c *gin.Context) {
	jobName := c.Query("job")
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job parameter"})
		return
	}

	if !handle.UserServices.CurrentPageDataExists(jobName) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or pages not initialized"})
		return
	}

	pages, exist := handle.UserServices.GetCurrentPageData(jobName)
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pages not found for the job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_pages": len(pages.PageItems),
		"pages":       pages.PageItems,
	})
}

// @Summary      Get all images for a job
// @Description  Returns all base64 encoded images for a specific page of a job
// @Tags         base64images
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Param        dataset  query  string  true  "Dataset name"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getImages [get]
func (handle *Handle) GetBase64ImageByPageIndex(c *gin.Context) {
	jobName := c.Query("job")
	datasetName := c.Query("dataset")
	pageIndex := c.Query("imageIndex")

	if jobName == "" || datasetName == "" || pageIndex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job, dataset or imageIndex parameter"})
		return
	}

	index, err := strconv.Atoi(pageIndex)
	if err != nil || index < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid imageIndex parameter"})
		return
	}

	if handle.UserServices.ImageCacheExists(jobName) {
		imagePaths, base64Images := handle.UserServices.GetBase64ImageCacheByPage(jobName, index)
		if len(imagePaths) == 0 || len(base64Images) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No images found for the specified page"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"base64_image": base64Images,
		})
	}

	imagePaths, base64Images := handle.UserServices.SetBase64ImageCacheByPage(jobName, index)
	if len(imagePaths) == 0 || len(base64Images) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found for the specified page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"base64_image": base64Images,
	})

}

// @Summary      Save pending review
// @Description  Save pending review data
// @Tags         review
// @Accept       json
// @Produce      json
// @Param        body  body  interface{}  true  "Pending review data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/savePendingReview [post]
func (handle *Handle) SavePendingReview(c *gin.Context) {
	var body interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	itemsLen := handle.UserServices.SavePendingReviewData(body)
	if itemsLen == 0 {
		log.Println("Failed to save pending review data")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  itemsLen,
	})
}

// @Summary      Get pending review
// @Description  Get all pending review data
// @Tags         review
// @Produce      json
// @Param        flatten  query  bool  false  "Flatten the result"
// @Success      200  {object}  interface{}
// @Router       /api/getPendingReview [get]
func (handle *Handle) GetPendingReview(c *gin.Context) {
	flatten := c.DefaultQuery("flatten", "false") == "true"

	if flatten {
		c.JSON(http.StatusOK, handle.UserServices.PendingReviewData)
		return
	}

	c.JSON(http.StatusOK, handle.UserServices.PendingReviewData)
}
