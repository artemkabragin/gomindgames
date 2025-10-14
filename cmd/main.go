package main

import (
	"log"
	"mindgames/internal/controller"
	"mindgames/internal/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize database before starting the server
	dsn := "host=localhost user=artembragin password=1337 dbname=auth_database port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	log.Println("Connected to database ", db.Table("users") != nil)

	if err := migrations.Run(db); err != nil {
		log.Fatalf("can't run migrations: %v", err)
	}

	_ = controller.NewController(controller.ControllerOptions{
		DB: db,
	})
}
