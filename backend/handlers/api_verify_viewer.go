// handlers/api_func_image.go
package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary      Get all job names
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
	type PageRequest struct {
		Job      string `json:"job" binding:"required"`
		PageSize int    `json:"pageSize" binding:"required,gt=0"`
	}

	var req PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: 'job' must be string, 'pageSize' must be a positive integer."})
		return
	}

	log.Printf("SetAllPageDetails called with job: %s and pageSize: %d", req.Job, req.PageSize)

	if !handle.UserServices.CurrentPageDataExists(req.Job) {
		handle.UserServices.CacheManager.ClearImageCacheStore(req.Job)
		handle.UserServices.ClearCurrentPageData()
		handle.UserServices.SetCurrentPageData(req.Job, req.PageSize)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All pages set successfully for job: " + req.Job,
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

// @Summary      Get base64image by page index
// @Description  Returns all base64 encoded images for a specific page of a job
// @Tags         images
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Param        dataset  query  string  true  "Dataset name"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getBase64ImageSet [get]
func (handle *Handle) GetBase64ImageByPageIndex(c *gin.Context) {
	jobName := c.Query("job")
	pageIndex := c.Query("pageIndex")

	if jobName == "" || pageIndex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job, dataset or imageIndex parameter"})
		return
	}

	index, err := strconv.Atoi(pageIndex)
	if err != nil || index < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid imageIndex parameter"})
		return
	}

	if handle.UserServices.ImageCacheExists(jobName, index) {
		imagePaths, base64Images := handle.UserServices.GetBase64ImageCacheByPage(jobName, index)
		if len(imagePaths) == 0 || len(base64Images) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No images found for the specified page"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"image_path":   imagePaths,
			"base64_image": base64Images,
		})
		return
	}

	imagePaths, base64Images := handle.UserServices.SetBase64ImageCacheByPage(jobName, index)
	log.Println("SetBase64ImageByPageIndex called with job:", jobName, "and pageIndex:", pageIndex)
	if len(imagePaths) == 0 || len(base64Images) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found for the specified page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"image_path":   imagePaths,
		"base64_image": base64Images,
	})

}

// @Summary      Get base64 image by image path
// @Description  Returns a base64 encoded image for a specific image path in a job
// @Tags         images
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Param        imagePath  query  string  true  "Image path"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getBase64Image [get]
func (Handle *Handle) GetBase64ImageByImagePath(c *gin.Context) {
	jobName := c.Query("job")
	imagePath := c.Query("imagePath")
	if jobName == "" || imagePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing imagePath parameter"})
		return
	}

	base64Image := Handle.UserServices.GetBase64ImageByPath(jobName, imagePath)
	if base64Image == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"base64_image": base64Image,
	})

}

// @Summary      Get image by page index
// @Description  Returns image names and paths for a specific page of a job
// @Tags         images
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Param        pageIndex  query  string  true  "Page index"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getImageSet [get]
func (handle *Handle) GetImageByPageIndex(c *gin.Context) {
	jobName := c.Query("job")
	pageIndex := c.Query("pageIndex")

	if jobName == "" || pageIndex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job or pageIndex parameter"})
		return
	}

	index, err := strconv.Atoi(pageIndex)
	if err != nil || index < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageIndex parameter"})
		return
	}

	if !handle.UserServices.CurrentPageDataExists(jobName) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image cache not found for the specified job"})
	}

	imageNames, imagePaths := handle.UserServices.GetImageCacheByPage(jobName, index)
	if len(imageNames) == 0 && len(imagePaths) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found for the specified page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image_name": imageNames,
		"image_path": imagePaths,
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

// @Summary      Get pending review image paths
// @Description  Returns a list of full image paths for all pending review items
// @Tags         review
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /api/getPendingReviewPaths [get]
func (handle *Handle) GetPendingReviewPaths(c *gin.Context) {
	items := handle.UserServices.GetPendingReviewItems()

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No pending review items found"})
		return
	}

	imagePaths := handle.UserServices.GetPendingReviewImagePaths()

	c.JSON(http.StatusOK, gin.H{
		"total_items": len(items),
		"image_paths": imagePaths,
	})
}

// @Summary      Get backup list
// @Description  Get list of all available backups
// @Tags         backup
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /api/getBackupList [get]
func (handle *Handle) GetBackupList(c *gin.Context) {
	backups, err := handle.UserServices.GetBackupList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get backup list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"backups": backups,
		"count":   len(backups),
	})
}

// @Summary      Restore from backup
// @Description  Restore pending review data from a specific backup
// @Tags         backup
// @Accept       json
// @Produce      json
// @Param        body  body  object{filename=string}  true  "Backup filename"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/restoreFromBackup [post]
func (handle *Handle) RestoreFromBackup(c *gin.Context) {
	var requestBody struct {
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err := handle.UserServices.RestoreFromBackup(requestBody.Filename)
	if err != nil {
		if err.Error() == "backup file not found: "+requestBody.Filename {
			c.JSON(http.StatusNotFound, gin.H{
				"error":    "Backup file not found",
				"filename": requestBody.Filename,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to restore from backup",
			"details": err.Error(),
		})
		return
	}

	items := handle.UserServices.GetPendingReviewItems()

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"message":        "Successfully restored from backup",
		"filename":       requestBody.Filename,
		"restored_items": len(items),
	})
}
