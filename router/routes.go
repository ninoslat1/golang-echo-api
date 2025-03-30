package router

import (
	authhandler "echo-api/handlers"
	repository "echo-api/repositories"
	service "echo-api/services"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(e *echo.Echo, log *logrus.Logger) {
	userRepo := repository.NewUserRepository(log)
	authService := service.NewAuthService(userRepo)

	authHandler := authhandler.NewAuthHandler(authService)

	e.POST("/login", authHandler.LoginHandler)
	e.POST("/register", authHandler.RegisterUserHandler)
	e.POST("/verify", authHandler.VerifyEmailHandler)
}
