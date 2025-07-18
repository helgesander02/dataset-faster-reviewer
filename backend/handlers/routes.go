package handlers

import (
	"github.com/gin-gonic/gin"
)

func (handle *Handle) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/folder-structure", handle.FolderStructureHandler)
		api.GET("/getJobs", handle.GetJobs)
		api.GET("/getDatasets", handle.GetDatasets)
		api.GET("/getImages", handle.GetImages)
		api.GET("/getBase64Images", handle.GetBase64Images)
		api.GET("/getAllPages", handle.GetAllPages)
		
		api.POST("/savePendingReview", handle.SavePendingReview)
		api.GET("/getPendingReview", handle.GetPendingReview)
		api.POST("/approvedRemove", handle.ApprovedRemove)
		api.POST("/unapprovedRemove", handle.UnApprovedRemove)
	}
}
