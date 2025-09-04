package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// Server
	Port string `env:"PORT" envDefault:"8080"`

	//// Google Cloud
	//ProjectID  string `env:"GCP_PROJECT_ID,required"`
	//BucketName string `env:"GCS_BUCKET_NAME,required"`
	//CredsPath  string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
	//CredsRaw   string `env:"GOOGLE_APPLICATION_CREDENTIALS_JSON"`
	//
	//// Logging
	//LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	//LogFormat string `env:"LOG_FORMAT" envDefault:"json"`
	//
	//// Derived
	//CredsJSON []byte `env:"-"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Defaults
	v.SetDefault("PORT", "8080")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	cfg := &Config{
		Port: v.GetString("PORT"),
		//ProjectID:  v.GetString("GCP_PROJECT_ID"),
		//BucketName: v.GetString("GCS_BUCKET_NAME"),
		//CredsPath:  v.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		//CredsRaw:   v.GetString("GOOGLE_APPLICATION_CREDENTIALS_JSON"),
		//LogLevel:   v.GetString("LOG_LEVEL"),
		//LogFormat:  v.GetString("LOG_FORMAT"),
	}

	// Materialize credentials JSON
	//if strings.TrimSpace(cfg.CredsRaw) != "" {
	//	cfg.CredsJSON = []byte(cfg.CredsRaw)
	//} else if p := strings.TrimSpace(cfg.CredsPath); p != "" {
	//	if b, err := os.ReadFile(p); err == nil {
	//		cfg.CredsJSON = b
	//	} else {
	//		return nil, err
	//	}
	//}

	return cfg, nil
}
