package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type appConfig struct {
	Environment string `envconfig:"ENV" default:"local"`
	Host        string `envconfig:"APP_HOST" default:"0.0.0.0"`
	Port        int    `envconfig:"APP_PORT" default:"3000"`
}

type mongoConfig struct {
	URI string `envconfig:"MONGO_URI"`
	DB  string `envconfig:"MONGO_DB" default:"test"`
}

type config struct {
	App        appConfig
	Mongo      mongoConfig
	JWT_Secret string `envconfig:"JWT_SECRET"`
}

var c config

func Load() {
	godotenv.Load()
	err := envconfig.Process("", &c)
	if err != nil {
		zap.L().Fatal("failed to load configuration", zap.Error(err))
	}
}

func Get() config {
	return c
}
