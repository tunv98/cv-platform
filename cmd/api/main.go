package main

import (
	"cv-platform/internal/adapter/http"
	"cv-platform/internal/config/dotenv"
	"cv-platform/pkg/healthcheck"
	"cv-platform/pkg/logger"
)

const (
	version = "1.0.0"
)

func main() {
	logger.Init("info", true) // Use console format for development
	log := logger.FLog()

	cfg, err := dotenv.Load()
	if err != nil {
		log.Errorf("failed to load config: %v", err)
		return
	}

	log.Infof("starting cv-platform API server: port=%s, version=%s", cfg.Port, version)

	r := http.NewRouter(cvUploadUC, profileStoreUC)
	healthcheck.Apply(r, version)

	log.Infof("server starting on address: %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Errorf("failed to run server: %v", err)
	}
}
