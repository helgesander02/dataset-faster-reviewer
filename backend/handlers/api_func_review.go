package handlers

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

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
	
	itemsLen := handle.DM.SavePendingReviewData(body)
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
		c.JSON(http.StatusOK, handle.DM.PendingReviewData)
		return
	}
	
	c.JSON(http.StatusOK, handle.DM.PendingReviewData)
}

// @Summary      Remove approved images
// @Description  Remove images that have been approved from pending review
// @Tags         review
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/approvedRemove [post]
func (handle *Handle) ApprovedRemove(c *gin.Context) {
	items := handle.DM.GetPendingReviewItems()
	var failedFiles []string
	
	for _, item := range items {
		filePath := handle.DM.ImageRoot + "/" + item.Job + "/" + item.Dataset + "/" + item.ImagePath
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove file: %s, error: %v", filePath, err)
			failedFiles = append(failedFiles, filePath)
		} else {
			log.Printf("Successfully removed file: %s", filePath)
		}
	}

	if len(failedFiles) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":       "partial_success",
			"failed_files": failedFiles,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// @Summary      Clear unapproved pending review data
// @Description  Clear all unapproved pending review data
// @Tags         review
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/unapprovedRemove [post]
func (handle *Handle) UnApprovedRemove(c *gin.Context) {
	handle.DM.ClearPendingReviewData()
	log.Println("PendingReviewData has been cleared.")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
