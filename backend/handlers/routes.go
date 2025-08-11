package handlers

import (
	"github.com/gin-gonic/gin"
)

func (handle *Handle) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/getJobs", handle.GetJobs)

		api.POST("/setAllPages", handle.SetAllPageDetails)
		api.GET("/getAllPages", handle.GetAllPageDetails)

		api.GET("/getBase64ImageSet", handle.GetBase64ImageByPageIndex)
		api.GET("/getBase64Image", handle.GetBase64ImageByImagePath)
		api.GET("/getImageSet", handle.GetImageByPageIndex)

		api.POST("/savePendingReview", handle.SavePendingReview)
		api.GET("/getPendingReview", handle.GetPendingReview)
		api.GET("/getPendingReviewPaths", handle.GetPendingReviewPaths)
	}
}
