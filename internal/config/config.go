package config

import (
	"os"
)

type DBConfig struct {
	User string
	Port string
	Name string
	Password string
	Host string
}

type Config struct {
	DB DBConfig
}

func New() *Config {
	return &Config{
		DB: DBConfig{
			User: getEnv("DB_USER", ""),
			Port: getEnv("DB_PORT", ""),
			Name: getEnv("DB_NAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host: getEnv("DB_HOST", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}

	return defaultVal
}