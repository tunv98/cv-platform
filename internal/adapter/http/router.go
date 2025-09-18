package http

import (
	"cv-platform/internal/adapter/http/handler"
	"cv-platform/internal/adapter/http/middleware"
	"cv-platform/internal/usecase"
	logger "cv-platform/pkg/log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(cvUC *usecase.CVUploadUC, profileUC *usecase.ProfileStoreUC) *gin.Engine {
	router := gin.New()

	// Recovery middleware (log panic ra stderr + logger)
	router.Use(gin.CustomRecoveryWithWriter(os.Stderr, func(c *gin.Context, err interface{}) {
		if e, ok := err.(string); ok {
			logger.FLog().Errorf("panic recovered: requestID = %s, path = %s, err = %s",
				c.GetString(middleware.RequestIDKey), c.Request.URL.Path, e)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	// Prometheus' metrics middleware
	router.Use(middleware.PrometheusMiddleware())
	// Logging middleware
	router.Use(middleware.RequestLogging([]string{"/health", "/info", "/metrics"}))
	cvHandler := handler.NewCVHandler(cvUC)
	profileHandler := handler.NewProfileHandler(profileUC)

	api := router.Group("/api/v1")
	cvApi := api.Group("/cvs")
	{
		cvApi.POST("/upload", cvHandler.StartUpload)
		cvApi.PUT("/:id", cvHandler.CompleteUpload)
	}
	profileApi := api.Group("/profiles")
	{
		profileApi.GET("/:id", profileHandler.GetProfile)
	}

	// Prometheus' metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
