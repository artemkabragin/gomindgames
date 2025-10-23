package main

import (
	"log"
	"mindgames/internal/controller"
	"mindgames/internal/kafka"
	"mindgames/internal/repository"
	"os"
	"strings"
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

	kafkaConfig := kafka.KafkaConfig{
		Brokers: strings.Split(getEnvOrDefault("KAFKA_BROKERS", "kafka:9091"), ","),
		Topic:   getEnvOrDefault("KAFKA_TOPIC", "uesr-events"),
		GroupID: getEnvOrDefault("KAFKA_GROUP_ID", "user-service"),
	}

	kafkaClient, err := kafka.NewKafkaClient(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %s", err.Error())
	}
	defer kafkaClient.Close()
	log.Println("Successfully connected to Kafka as producer")

	controller := controller.NewController(controller.ControllerOptions{
		DB:          db,
		KafkaClient: kafkaClient,
	})

	if controller == nil {
		log.Fatal("Failed to start controller")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
