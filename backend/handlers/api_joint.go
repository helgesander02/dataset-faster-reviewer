package handlers

import (
	"backend/src/services"
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

	itemsLen := handle.JointServices.SavePendingReviewData(body)
	if itemsLen == 0 {
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
// @Success      200  {object}  interface{}
// @Router       /api/getPendingReview [get]
func (handle *Handle) GetPendingReview(c *gin.Context) {
	c.JSON(http.StatusOK, handle.JointServices.PendingReviewData)
}

// @Summary      Clear pending review
// @Description  Clear all pending review data and caches
// @Tags         review
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/clearPendingReview [post]
func (handle *Handle) ClearPendingReview(c *gin.Context) {
	handle.JointServices.ClearPendingReview()
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "All pending review data cleared",
	})
}

// @Summary      Get review image
// @Description  Get original image for review item by job, dataset, and image name (direct file serving)
// @Tags         review
// @Produce      image/jpeg,image/png
// @Param        job         query    string  true  "Job name"
// @Param        dataset     query    string  true  "Dataset name"
// @Param        imageName   query    string  true  "Image name"
// @Success      200  {file}  binary
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/getReviewImage [get]
func (handle *Handle) GetReviewImage(c *gin.Context) {
	job := c.Query("job")
	dataset := c.Query("dataset")
	imageName := c.Query("imageName")

	if job == "" || dataset == "" || imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: job, dataset, imageName",
		})
		return
	}

	// Construct the full image path using GetImageRoot for consistency
	root := services.GetImageRoot()
	imagePath := root + "/" + job + "/" + dataset + "/image/" + imageName

	// Serve the file directly as binary
	c.File(imagePath)
}

// @Summary      Get pending review image paths
// @Description  Returns a list of full image paths for all pending review items
// @Tags         review
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /api/getPendingReviewPaths [get]
func (handle *Handle) GetPendingReviewPaths(c *gin.Context) {
	items := handle.JointServices.GetPendingReviewItems()

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No pending review items found"})
		return
	}

	imagePaths := handle.JointServices.GetPendingReviewImagePaths()

	c.JSON(http.StatusOK, gin.H{
		"total_items": len(items),
		"image_paths": imagePaths,
	})
}

// @Summary      Get paginated review images
// @Description  Returns paginated review images with base64 data for review modal
// @Tags         review
// @Produce      json
// @Param        page   query    int  false  "Page number (0-indexed)"  default(0)
// @Param        limit  query    int  false  "Images per page"          default(9)
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /api/getPendingReviewImages [get]
func (handle *Handle) GetPendingReviewImages(c *gin.Context) {
	// Parse pagination parameters
	pageStr := c.DefaultQuery("page", "0")
	limitStr := c.DefaultQuery("limit", "9")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page parameter",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	// Get all review items
	allItems := handle.JointServices.GetPendingReviewItems()

	if len(allItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No pending review items found"})
		return
	}

	// Calculate pagination bounds
	start := page * limit
	if start >= len(allItems) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Page number out of range",
		})
		return
	}

	end := start + limit
	if end > len(allItems) {
		end = len(allItems)
	}

	// Get items for this page
	pageItems := allItems[start:end]

	// Extract image paths from items
	imagePaths := make([]string, len(pageItems))
	for i, item := range pageItems {
		imagePaths[i] = item.ImagePath
	}

	// Get base64 images using Review cache (separate from ImageGrid cache)
	base64Images := handle.UserServices.GetOrCreateReviewBase64Images(imagePaths)

	c.JSON(http.StatusOK, gin.H{
		"image_path":   imagePaths,
		"base64_image": base64Images,
		"total_items":  len(allItems),
		"page":         page,
		"limit":        limit,
		"page_items":   len(pageItems),
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
	backups, err := handle.JointServices.GetBackupList()
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

	err := handle.JointServices.RestoreFromBackup(requestBody.Filename)
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

	items := handle.JointServices.GetPendingReviewItems()

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"message":        "Successfully restored from backup",
		"filename":       requestBody.Filename,
		"restored_items": len(items),
	})
}

// @Summary      Delete selected images
// @Description  Delete selected images from pending review and clear caches
// @Tags         review
// @Accept       json
// @Produce      json
// @Param        body  body  interface{}  true  "Array of images to delete"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/deleteSelectedImages [post]
func (handle *Handle) DeleteSelectedImages(c *gin.Context) {
	var body interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Delete physical files and get result
	result, err := handle.JointServices.DeleteSelectedImages(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	// Clean up caches for deleted images (non-blocking, errors are logged but don't fail the request)
	handle.JointServices.CleanupDeletedImagesFromCache(handle.UserServices, result)

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"deleted_count": result.DeletedCount,
		"cache_cleared": result.CacheCleared,
		"affected_jobs": result.AffectedJobs,
	})
}
