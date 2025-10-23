package controller

import (
	"mindgames/internal/handler"
	"mindgames/internal/kafka"
	"mindgames/internal/repository"
	"mindgames/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type ControllerOptions struct {
	DB          *gorm.DB
	KafkaClient kafka.IKafkaClient
}

type Controller struct {
	db          *gorm.DB
	kafkaClient kafka.IKafkaClient
	userService service.UserService
	userHandler handler.UserHandler
}

func NewController(opts ControllerOptions) *Controller {
	userRepo := repository.NewUserRepository(opts.DB)
	tokenRepo := repository.NewTokenRepository(opts.DB)

	eventProducer := kafka.NewEventProducer(opts.KafkaClient)

	userService := service.NewUserService(userRepo, eventProducer)
	tokenService := service.NewTokenService(tokenRepo)

	userHandler := handler.NewUserHandler(userService, tokenService)
	testHandler := handler.NewTestHandler()

	e := initEcho()

	e.POST("/auth/register", userHandler.PostRegister)
	e.POST("/auth/login", userHandler.PostLogin)
	e.POST("/auth/refresh", userHandler.PostRefresh)
	e.GET("/getTest", testHandler.GetTest)

	err := e.Start("0.0.0.0:8081")
	if err != nil {
		return nil
	}

	return &Controller{
		db:          opts.DB,
		kafkaClient: opts.KafkaClient,
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
