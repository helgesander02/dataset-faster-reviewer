package handlers

import (
	"github.com/gin-gonic/gin"
)

func (handle *Handle) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/getJobs", handle.GetJobs)
		api.GET("/setAllPages", handle.SetAllPageDetails)
		api.GET("/getAllPages", handle.GetAllPageDetails)
		api.GET("/getBase64Images", handle.GetBase64ImageByPageIndex)

		api.POST("/savePendingReview", handle.SavePendingReview)
		api.GET("/getPendingReview", handle.GetPendingReview)
	}
}
