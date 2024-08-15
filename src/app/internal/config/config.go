package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	SessionCookieName string
	Dsn               string
}

func setDsn() string {
	// TODO: setup for mysql
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s database=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)
	return dsn
}

func loadConfig() (*Config, error) {

	var cfg Config
	dsn := setDsn()
	cfg.Dsn = dsn
	cfg.Port = "4000"
	cfg.SessionCookieName = "session"

	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
