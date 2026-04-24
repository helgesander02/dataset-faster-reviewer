package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
		Job          string `json:"job" binding:"required"`
		ImagePerPage int    `json:"image_per_page" binding:"required,gt=0"`
	}

	var req PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: 'job' must be string, 'image_per_page' must be a positive integer",
		})
		return
	}

	log.Printf("[SetAllPageDetails] REQUEST - job: %s, image_per_page: %d", req.Job, req.ImagePerPage)

	pageDataExists := handle.UserServices.CurrentPageDataExists(req.Job)
	log.Printf("[SetAllPageDetails] CurrentPageDataExists(%s): %v", req.Job, pageDataExists)

	if !pageDataExists {
		log.Printf("[SetAllPageDetails] Clearing and setting new page data for job: %s", req.Job)
		handle.UserServices.ClearImageCache(req.Job)
		handle.UserServices.ClearCurrentPageData()
		handle.UserServices.SetCurrentPageData(req.Job, req.ImagePerPage)
	} else {
		log.Printf("[SetAllPageDetails] Page data already exists for job: %s, skipping setup", req.Job)
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
	log.Printf("[GetAllPageDetails] REQUEST - job: %s", jobName)

	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job parameter"})
		return
	}

	pageDataExists := handle.UserServices.CurrentPageDataExists(jobName)
	log.Printf("[GetAllPageDetails] CurrentPageDataExists(%s): %v", jobName, pageDataExists)

	if !pageDataExists {
		log.Printf("[GetAllPageDetails] ERROR - Job not found or pages not initialized for: %s", jobName)
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or pages not initialized"})
		return
	}

	pages, exist := handle.UserServices.GetCurrentPageData(jobName)
	if !exist {
		log.Printf("[GetAllPageDetails] ERROR - Pages not found after existence check for job: %s", jobName)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pages not found for the job"})
		return
	}

	// Use PageItemsReadOnly for performance - avoid copying large slice before JSON serialization
	pageItems := pages.PageItemsReadOnly()
	log.Printf("[GetAllPageDetails] SUCCESS - Returning %d pages for job: %s", len(pageItems), jobName)
	c.JSON(http.StatusOK, gin.H{
		"total_pages": len(pageItems),
		"pages":       pageItems,
	})
}

// @Summary      Get job metadata (lightweight)
// @Description  Returns lightweight job metadata including per-page dataset mapping and total pages without image details
// @Tags         pages
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/getJobMetadata [get]
func (handle *Handle) GetJobMetadata(c *gin.Context) {
	jobName := c.Query("job")
	log.Printf("[GetJobMetadata] REQUEST - job: %s", jobName)

	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job parameter"})
		return
	}

	pageDataExists := handle.UserServices.CurrentPageDataExists(jobName)
	log.Printf("[GetJobMetadata] CurrentPageDataExists(%s): %v", jobName, pageDataExists)

	if !pageDataExists {
		log.Printf("[GetJobMetadata] ERROR - Job not found or pages not initialized for: %s", jobName)
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or pages not initialized"})
		return
	}

	pages, exist := handle.UserServices.GetCurrentPageData(jobName)
	if !exist {
		log.Printf("[GetJobMetadata] ERROR - Pages not found after existence check for job: %s", jobName)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pages not found for the job"})
		return
	}

	// Use cached datasets - O(1) operation!
	datasetNames := pages.GetDatasetNames()
	totalPages := pages.Len()

	log.Printf("[GetJobMetadata] SUCCESS - job: %s, total_pages: %d, datasets: %d", jobName, totalPages, len(datasetNames))
	c.JSON(http.StatusOK, gin.H{
		"job_name":      jobName,
		"total_pages":   totalPages,
		"dataset_names": datasetNames,
	})
}

// @Summary      Get page by page index
// @Description  Returns a specific page by its index for a given job
// @Tags         pages
// @Produce      json
// @Param        job  query  string  true  "Job name"
// @Param        pageIndex  query  string  true  "Page index"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/getPageByPageIndex [get]
func (handle *Handle) GetPageByPageIndex(c *gin.Context) {
	jobName := c.Query("job")
	pageIndex := c.Query("pageIndex")

	if jobName == "" || pageIndex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job or pageIndex parameter"})
		return
	}

	index, err := parsePageIndex(pageIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageIndex parameter"})
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

	pageItem, found := pages.PageAt(index)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page index out of range"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page_index": pageIndex,
		"page":       pageItem,
	})
}

func parsePageIndex(pageIndex string) (int, error) {
	index, err := strconv.Atoi(pageIndex)
	if err != nil || index < 0 {
		return 0, err
	}
	return index, nil
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

	log.Printf("[GetBase64ImageByPageIndex] REQUEST - job: %s, pageIndex: %s", jobName, pageIndex)

	if jobName == "" || pageIndex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job or pageIndex parameter"})
		return
	}

	index, err := parsePageIndex(pageIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageIndex parameter"})
		return
	}

	imagePaths, base64Images := handle.getOrCreateBase64ImageCache(jobName, index, pageIndex)

	// Check if task was cancelled (returns nil, nil)
	if imagePaths == nil || base64Images == nil {
		log.Printf("[GetBase64ImageByPageIndex] ERROR - Task was cancelled for job: %s, page: %d", jobName, index)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Image processing was cancelled, please retry",
			"code":  "TASK_CANCELLED",
		})
		return
	}

	// Check for empty arrays
	if len(imagePaths) == 0 || len(base64Images) == 0 {
		log.Printf("[GetBase64ImageByPageIndex] ERROR - No images found for job: %s, page: %d", jobName, index)
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found for the specified page"})
		return
	}

	// Check for empty base64 strings in response
	emptyCount := 0
	for _, b64 := range base64Images {
		if b64 == "" {
			emptyCount++
		}
	}
	if emptyCount > 0 {
		log.Printf("[GetBase64ImageByPageIndex] WARNING - Returning %d empty base64 strings out of %d total for job: %s, page: %d", emptyCount, len(base64Images), jobName, index)
	} else {
		log.Printf("[GetBase64ImageByPageIndex] SUCCESS - Returning %d images for job: %s, page: %d", len(base64Images), jobName, index)
	}

	c.JSON(http.StatusOK, gin.H{
		"image_path":   imagePaths,
		"base64_image": base64Images,
	})
}

func (handle *Handle) getOrCreateBase64ImageCache(jobName string, index int, pageIndex string) ([]string, []string) {
	if handle.UserServices.ImageCacheExists(jobName, index) {
		log.Printf("[getOrCreateBase64ImageCache] Cache HIT for job: %s, page: %d", jobName, index)
		imagePaths, base64Images := handle.UserServices.GetBase64ImageCacheByPage(jobName, index)
		log.Printf("[getOrCreateBase64ImageCache] Retrieved from cache: %d paths, %d images", len(imagePaths), len(base64Images))
		return imagePaths, base64Images
	}

	log.Printf("[getOrCreateBase64ImageCache] Cache MISS - Creating base64 image cache for job: %s, page: %s", jobName, pageIndex)
	return handle.UserServices.SetBase64ImageCacheByPage(jobName, index)
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
		return
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
