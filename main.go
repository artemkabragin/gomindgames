package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Initialize database before starting the server
	dsn := "host=localhost user=artembragin password=1337 dbname=auth_database port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	log.Println("Connected to database ", db.Table("users") != nil)

	if err := db.AutoMigrate(&User{}, &RefreshToken{}); err != nil {
		log.Fatalf("can't migrate database: %v", err)
	}

	// Start Echo after DB is ready
	e := initEcho()
	e.Start("localhost:8080")
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/auth/register", postRegister)
	e.POST("/auth/login", postLogin)
	e.POST("/auth/refresh", postRefresh)

	return e
}
