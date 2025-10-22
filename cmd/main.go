package main

import (
	"log"
	"mindgames/internal/controller"
	"mindgames/internal/repository"
	"os"
	"sync"
)

func main() {
	log.Println("Starting application...")
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

	controller := controller.NewController(controller.ControllerOptions{
		DB: db,
	})

	if controller == nil {
		log.Fatal("Failed to start controller")
	}

	// Keep the main goroutine alive
	select {}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
