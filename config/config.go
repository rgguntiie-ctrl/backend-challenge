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

type PsqlConfig struct {
	Host string `envconfig:"PSQL_HOST" default:"localhost"`
	DB   string `envconfig:"PSQL_DB" default:"test"`
	User string `envconfig:"PSQL_USER" default:"postgres"`
	Pass string `envconfig:"PSQL_PASS" default:"password"`
	Port string `envconfig:"PSQL_PORT" default:"5432"`
}

type redisConfig struct {
	Addr     string `envconfig:"REDIS_ADDRESS" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

type config struct {
	App        appConfig
	Mongo      mongoConfig
	JWT_Secret string `envconfig:"JWT_SECRET"`
	Psql       PsqlConfig
	Redis      redisConfig
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
