package controller

import (
	"mindgames/internal/handler"
	"mindgames/internal/repository"
	"mindgames/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type ControllerOptions struct {
	DB *gorm.DB
}

type Controller struct {
	db          *gorm.DB
	userService service.IUserService
	userHandler handler.IUserHandler
}

func NewController(opts ControllerOptions) *Controller {
	userRepo := repository.UserRepo(opts.DB)
	tokenRepo := repository.NewTokenRepository(opts.DB)

	userService := service.UserService(userRepo)
	tokenService := service.TokenService(tokenRepo)

	userHandler := handler.NewUserHandler(userService, tokenService)
	testHandler := handler.TestHandler()

	e := initEcho()

	e.POST("/auth/register", userHandler.PostRegister)
	e.POST("/auth/login", userHandler.PostLogin)
	e.POST("/auth/refresh", userHandler.PostRefresh)
	e.GET("/getTest", testHandler.GetTest)

	err := e.Start("localhost:8081")
	if err != nil {
		return nil
	}

	return &Controller{
		db:          opts.DB,
		userService: userService,
		userHandler: userHandler,
	}
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	return e
}
