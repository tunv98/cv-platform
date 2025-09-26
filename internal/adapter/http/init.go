package http

import (
	config "cv-platform/internal/config/ymlreader"
	"cv-platform/internal/usecase"
)

func Init() {
	cfg := config.GetConfig()
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

}
