package http

import (
	"cv-platform/internal/adapter/http/handler"
	"cv-platform/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewRouter(cvUC *usecase.CVUploadUC, profileUC *usecase.ProfileStoreUC) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	cvApi := api.Group("/cvs")
	{
		cvApi.POST("/upload", handler.NewCVHandler(cvUC).StartUpload)
		cvApi.PUT("/:id", handler.NewCVHandler(cvUC).CompleteUpload)
	}
	profileApi := api.Group("/profiles")
	{
		profileApi.GET("/:id", handler.NewProfileHandler(profileUC).GetProfile)
	}
	return router
}
