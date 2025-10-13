package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	initEcho()
}

func initEcho() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/auth/register", postRegister)
	e.POST("/auth/login", postLogin)
	e.POST("/auth/refresh", postRefresh)

	e.Start("localhost:8080")
}
