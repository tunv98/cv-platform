package http

import (
	"cv-platform/internal/adapter/http/handler"
	"cv-platform/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewRouter(uc *usecase.CVUploadUC) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/cvs/upload", handler.NewCVHandler(uc).StartUpload)
		api.PUT("/cvs/:id", handler.NewCVHandler(uc).CompleteUpload)
	}
	return router
}
