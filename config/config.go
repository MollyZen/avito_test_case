package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type (
	// Config -.
	Config struct {
		App      `yaml:"main"`
		HTTP     `yaml:"http"`
		Log      `yaml:"log"`
		Postgres `yaml:"database"`
	}

	// App -.
	App struct {
		Name string `yaml:"name" env:"TC_APP_NAME" env-default:"AVITO_TEST_CASE_SALTYKOV"`
	}

	// HTTP -.
	HTTP struct {
		Port int `yaml:"port" env:"TC_APP_VERSION" validate:"gt=0" env-default:"80"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"level" env:"TC_LOG_LEVEL" env-default:"WARN"`
	}

	// Postgres -.
	Postgres struct {
		Host              string `yaml:"host" env:"TC_PG_HOST" env-default:"localhost" validate:"required"`
		Port              int    `yaml:"port" env:"TC_PG_PORT" env-default:"5432" validate:"required"`
		User              string `yaml:"user" env:"TC_PG_USER" env-default:"postgres" validate:"required"`
		Password          string `yaml:"password" env:"TC_PG_PASSWORD" env-default:"postgres" validate:"required"`
		DB                string `yaml:"db" env:"TC_PG_DB" env-default:"avito_test_case" validate:"required"`
		PoolMaxOpen       int    `yaml:"poolMaxOpen" env:"TC_PG_POOL_MAX_OPEN" env-default:"10" validate:"required"`
		PoolMaxIdle       int    `yaml:"poolMaxIdle" env:"TC_PG_POOL_MAX_IDLE" env-default:"10" validate:"required"`
		PoolMaxLifetime   int    `yaml:"poolMaxLifetime" env:"TC_PG_POOL_MAX_LIFETIME" env-default:"3" validate:"required"`
		ReconnectAttempts int    `yaml:"ReconnectAttempts" env:"TC_PG_RECONNECT_ATTEMPTS" env-default:"3" validate:"required"`
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

	log.Print(cfg)

	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
