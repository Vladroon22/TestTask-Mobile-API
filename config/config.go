package config

import "os"

type Config struct {
	DB  string
	JWT string
}

func CreateConfig() *Config {
	return &Config{
		DB:  getEnv("DB", ""),
		JWT: getEnv("KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
