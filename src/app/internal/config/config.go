package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	ServerAddr        string
	SessionCookieName string
	Dsn               string
}

func setDsn() string {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("DB_ADMIN_USER")
	password := os.Getenv("DB_ADMIN_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		return ""
	}

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
	cfg.SessionCookieName = os.Getenv("SESSION_COOKIE_NAME")
	cfg.ServerAddr = fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	if dsn == "" {

		return nil, errors.New("No DSN found.")
	}

	if cfg.SessionCookieName == "" {
		return nil, errors.New("No session cookie name set.")
	}

	if cfg.ServerAddr == "" {
		return nil, errors.New("No server address set.")
	}

	cfg.Dsn = dsn

	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
