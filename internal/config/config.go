package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

var ErrConfigInit = errors.New("failed to init config")

// Config struct that depends configuration of App.
type Config struct {
	Server struct {
		Port         string `yaml:"port" env:"PORT" env-default:"8080"`
		WriteTimeout int    `yaml:"write_timeout" env-default:"15"`
		ReadTimeout  int    `yaml:"read_timeout" env-default:"15"`
		IdleTimeout  int    `yaml:"idle_timeout" env-default:"30"`
	}
	Postgres struct {
		Host     string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"PG_PORT" env-default:"5432"`
		User     string `yaml:"user" env:"PG_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"PG_PASSWORD" env-default:"postgres"`
		DBName   string `yaml:"db_name" env:"PG_NAME" env-default:"finance"`
	}
	Logger struct {
		Path  string `yaml:"path" env:"LOG_PATH" env-default:"./logs/logs.log"`
		Level int    `yaml:"level" env:"LOG_LEVEL" env-default:"6"`
	}
}

// GetConfig return pointer to config. Config is singleton.
func GetConfig(path string) (c Config, err error) {
	log.Print("reading server config file")
	if path == "" {
		path = "./config/config.yaml"
	}

	instance := Config{}
	if err = cleanenv.ReadConfig(path, &instance); err != nil {
		return Config{}, fmt.Errorf("%s: %w", ErrConfigInit.Error(), err)
	}

	return instance, nil
}
