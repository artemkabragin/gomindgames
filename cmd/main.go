package main

import (
	"context"
	"log"
	"mindgames/internal/controller"
	"mindgames/internal/kafka"
	"mindgames/internal/repository"
	"os"
	"strings"
)

func main() {
	log.Println("Starting application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	brokers := strings.Split(getEnvOrDefault("KAFKA_BROKERS", "kafka:9091"), ",")
	kafkaConfig := kafka.KafkaConfig{
		Brokers: brokers,
		Topic:   getEnvOrDefault("KAFKA_TOPIC", "user-events"),
		GroupID: getEnvOrDefault("KAFKA_GROUP_ID", "user-service"),
	}

	kafkaCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	kafkaClient, err := kafka.NewKafkaClient(kafkaCtx, kafkaConfig)
	if err != nil {
		log.Fatalf("failed to initialize Kafka: %s", err.Error())
	}
	defer kafkaClient.Close()
	log.Println("Successfully connected to Kafka as producer")

	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	controllerOptions := controller.ControllerOptions{
		DB:          db,
		KafkaClient: kafkaClient,
	}

	controller := controller.NewController(
		consumerCtx,
		controllerOptions,
	)

	if controller == nil {
		log.Fatal("failed to start controller")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
