package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	appConfig *AppConfig
	lock      = &sync.Mutex{}
)

type AppConfig struct {
	App struct {
		Name       string   `toml:"name"`
		Env        string   `toml:"env"`
		Port       uint16   `toml:"port"`
		Address    string   `toml:"address"`
		CorsOrigin []string `toml:"cors_origin"`
	} `toml:"app"`
	Database struct {
		Name                  string `toml:"name"`
		Username              string `toml:"username"`
		Password              string `toml:"password"`
		Address               string `toml:"address"`
		Port                  uint16 `toml:"port"`
		MaxConnection         int32  `toml:"maxConnection"`
		MinConnection         int32  `toml:"minConnection"`
		MaxConnectionLifeTime uint16 `toml:"maxConnectionLifeTime"`
		MinConnectionLifeTime uint16 `toml:"minConnectionLifeTime"`
		HealthCheckPeriod     string `toml:"healthCheckPeriod"`
	} `toml:"database"`
}

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		var err error
		appConfig, err = loadConfig()
		if err != nil {
			log.Fatal("Cant load config: ", err)
		}
	}

	return appConfig
}

func loadConfig() (*AppConfig, error) {
	viper.AddConfigPath("./config/")
	viper.SetConfigType("toml")

	env := os.Getenv("APP_ENV")
	if env == "PRODUCTION" || env == "DEVELOPMENT" {
		viper.SetConfigName("app")
	} else {
		viper.SetConfigName("test")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, use environment variable APP_ENV")
			return nil, err
		}
	}

	var config AppConfig
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &config, nil
}
