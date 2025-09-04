package main

import (
	"cv-platform/internal/adapter/http"
	"cv-platform/internal/config"
	logger "cv-platform/internal/log"
	"cv-platform/internal/usecase"
)

func main() {
	logger.Init("info", false) // Use console format for development
	log := logger.Simple()

	cfg, err := config.Load()
	if err != nil {
		log.Errorf("failed to load config: %v", err)
		return
	}

	log.Infof("starting cv-platform API server: port=%s, version=%s", cfg.Port, "1.0.0")

	// ctx := context.Background()
	// storage, err := gcp.NewGCSStorage(ctx, cfg.BucketName, cfg.CredsJSON)
	// if err != nil {
	// 	log.Errorf("failed to create gcs storage: %v", err)
	// 	return
	// }

	// repo, err := gcp.NewFirestoreCVRepo(ctx, cfg.ProjectID, cfg.CredsJSON)
	// if err != nil {
	// 	log.Errorf("failed to create firestore cv repo: %v", err)
	// 	return
	// }

	// cvUploadUC := usecase.NewCVUploadUC(storage, repo)
	var cvUploadUC *usecase.CVUploadUC
	profileStoreUC := usecase.NewProfileStoreUC()

	r := http.NewRouter(cvUploadUC, profileStoreUC)

	log.Infof("server starting on address: :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Errorf("failed to run server: %v", err)
	}
}
