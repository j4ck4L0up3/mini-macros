package config

import "github.com/kelseyhightower/envconfig"

type DsnConfig struct {
	Host         string `envconfig:"GO_DB_HOST"     default:"localhost"`
	User         string `envconfig:"GO_DB_USER"     default:"postgres"`
	Password     string `envconfig:"GO_DB_PASSWORD" default:"postgres"`
	DatabaseName string `envconfig:"GO_DB_NAME"     default:"test"`
	Port         string `envconfig:"GO_DB_PORT"     default:"5432"`
}

type Config struct {
	Port              string `envconfig:"PORT"                default:"localhost:4000"`
	SessionCookieName string `envconfig:"SESSION_COOKIE_NAME" default:"session"`
	DsnConfig         DsnConfig
}

func loadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
