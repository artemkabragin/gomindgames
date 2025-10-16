package main

import (
	"mindgames/internal/controller"
	"mindgames/internal/repository"
	"os"
)

func main() {
	var db = repository.NewGormDB(
		repository.Config{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			Username: getEnvOrDefault("DB_USER", "artembragin"),
			Password: getEnvOrDefault("DB_PASSWORD", "1337"),
			DBName:   getEnvOrDefault("DB_NAME", "auth_database"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
	)
	_ = controller.NewController(controller.ControllerOptions{
		DB: db,
	})
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
