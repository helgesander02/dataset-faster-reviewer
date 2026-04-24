package handlers

import (
	"backend/src/services"
	"context"

	"github.com/gin-gonic/gin"
)

type Handle struct {
	UserServices  *services.UserServices
	JointServices *services.JointServices
	ctx           context.Context
}

func NewHandle(ctx context.Context, root string, backupDir string) *Handle {
	services.SetConfig(root, backupDir)
	js := services.NewJointServices(ctx)
	us := services.NewUserServices()

	services.CheckServicesState(us, js)
	return &Handle{
		UserServices:  us,
		JointServices: js,
		ctx:           ctx,
	}
}

func (handle *Handle) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/getJobs", handle.GetJobs)

		api.POST("/setAllPages", handle.SetAllPageDetails)
		api.GET("/getAllPages", handle.GetAllPageDetails)
		api.GET("/getJobMetadata", handle.GetJobMetadata)
		api.GET("/getPage", handle.GetPageByPageIndex)

		api.GET("/getBase64ImageSet", handle.GetBase64ImageByPageIndex)
		api.GET("/getBase64Image", handle.GetBase64ImageByImagePath)
		api.GET("/getImageSet", handle.GetImageByPageIndex)

		api.POST("/savePendingReview", handle.SavePendingReview)
		api.GET("/getPendingReview", handle.GetPendingReview)
		api.POST("/clearPendingReview", handle.ClearPendingReview)
		api.GET("/getReviewImage", handle.GetReviewImage)
		api.GET("/getPendingReviewPaths", handle.GetPendingReviewPaths)
		api.GET("/getPendingReviewImages", handle.GetPendingReviewImages)
		api.POST("/deleteSelectedImages", handle.DeleteSelectedImages)

		api.GET("/getBackupList", handle.GetBackupList)
		api.POST("/restoreFromBackup", handle.RestoreFromBackup)
	}
}
