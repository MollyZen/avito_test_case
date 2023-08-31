package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	Config struct {
		App      `yaml:"main"`
		HTTP     `yaml:"http"`
		Log      `yaml:"log"`
		Postgres `yaml:"database"`
	}

	App struct {
		Name                       string `yaml:"name" env:"TC_APP_NAME" env-default:"AVITO_TEST_CASE_SALTYKOV"`
		AssignmentCleaningInterval int64  `yaml:"assignmentCleaningInterval" env:"TC_APP_ASSIGNMENT_CLEANING_INTERVAL" env-default:"3600" validate:"required"`
	}

	HTTP struct {
		Port int `yaml:"port" env:"TC_HTTP_PORT" validate:"gt=0" env-default:"8080"`
	}

	Log struct {
		Level string `yaml:"level" env:"TC_LOG_LEVEL" env-default:"INFO"`
	}

	Postgres struct {
		Host              string `yaml:"host" env:"TC_PG_HOST" env-default:"localhost" validate:"required"`
		Port              int32  `yaml:"port" env:"TC_PG_PORT" env-default:"8001" validate:"gt=0,required"`
		User              string `yaml:"user" env:"TC_PG_USER" env-default:"postgres" validate:"required"`
		Password          string `yaml:"password" env:"TC_PG_PASSWORD" env-default:"postgres" validate:"required"`
		DB                string `yaml:"db" env:"TC_PG_DB" env-default:"avito_test_case" validate:"required"`
		PoolMaxOpen       int32  `yaml:"poolMaxOpen" env:"TC_PG_POOL_MAX_OPEN" env-default:"10" validate:"required"`
		PoolMaxIdle       int32  `yaml:"poolMaxIdle" env:"TC_PG_POOL_MAX_IDLE" env-default:"10" validate:"required"`
		PoolMaxLifetime   int32  `yaml:"poolMaxLifetime" env:"TC_PG_POOL_MAX_LIFETIME" env-default:"600" validate:"required"`
		ReconnectAttempts int32  `yaml:"reconnectAttempts" env:"TC_PG_RECONNECT_ATTEMPTS" env-default:"3"`
		LogLevel          string `yaml:"logLevel" env:"TC_PG_LOG_LEVEL" env-default:"WARN" validate:"required"`
	}
)

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	//todo: добавить считывание из командой строки
	if _, err := os.Stat("config.yml"); err == nil {
		if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
			return nil, err
		}
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
