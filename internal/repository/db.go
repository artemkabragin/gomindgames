package repository

import (
	"fmt"
	"log"
	"mindgames/internal/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewGormDB(cfg Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user= %s password= %s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	log.Println("Connected to database ", db.Table("users") != nil)

	if err := migrations.Run(db); err != nil {
		log.Fatalf("can't run migrations: %v", err)
	}
	return db
}
