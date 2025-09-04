package main

import (
	"cv-platform/internal/adapter/http"
	"cv-platform/internal/config"
	logger "cv-platform/internal/log"
	"cv-platform/internal/usecase"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.L().Fatal("failed to load config", zap.Error(err))
	}

	// ctx := context.Background()
	// storage, err := gcp.NewGCSStorage(ctx, cfg.BucketName, cfg.CredsJSON)
	// if err != nil {
	// 	logger.L().Fatal("failed to create gcs storage", zap.Error(err))
	// }

	// repo, err := gcp.NewFirestoreCVRepo(ctx, cfg.ProjectID, cfg.CredsJSON)
	// if err != nil {
	// 	logger.L().Fatal("failed to create firestore cv repo", zap.Error(err))
	// }

	// cvUploadUC := usecase.NewCVUploadUC(storage, repo)
	var cvUploadUC *usecase.CVUploadUC
	profileStoreUC := usecase.NewProfileStoreUC()

	r := http.NewRouter(cvUploadUC, profileStoreUC)
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.L().Fatal("failed to run server", zap.Error(err))
	}
}
