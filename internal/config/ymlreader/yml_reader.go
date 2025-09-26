package ymlreader

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Service *Service `mapstructure:"service"`
}

type Service struct {
	Name    string  `mapstructure:"name"`
	Version string  `mapstructure:"version"`
	HTTP    *Listen `mapstructure:"http"`
	GRPC    *Listen `mapstructure:"grpc"`
	Env     string  `mapstructure:"env"`
}

type Listen struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

const (
	defaultPort = 8080
	defaultHost = "0.0.0.0"
)

var (
	defaultConfig []byte
	cfg           = &Config{}
)

func GetConfig() *Config {
	return cfg
}

func init() {
	viper.AutomaticEnv()
	configFile := viper.GetString("CONFIG_PATH")
	viper.SetConfigType("yaml")
	viper.AllowEmptyEnv(false)
	if configFile == "" {
		if err := viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
			panic(err)
		}
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
}

func (l *Listen) Address() string {
	if l == nil {
		return fmt.Sprintf("%s:%d", defaultHost, defaultPort)
	}
	if l.Host == "" {
		l.Host = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%s", l.Host, l.Port)
}
