package router

import (
	"echo-api/handlers"
	"echo-api/middlewares"
	"echo-api/repositories"
	"echo-api/services"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(e *echo.Echo, log *logrus.Logger) {
	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)

	authHandler := handlers.NewAuthHandler(authService, log)

	e.GET("/", handlers.LoginPageHandler)
	e.GET("/home", handlers.HomePageHandler)

	e.POST("/login", authHandler.LoginHandler)
	e.POST("/register", authHandler.RegisterUserHandler, middlewares.CookiePageMiddleware)
	e.POST("/verify", authHandler.VerifyEmailHandler)
	e.POST("/deactivate", authHandler.SoftDeleteUserHandler, middlewares.CookieAuthMiddleware)
	e.POST("/delete", authHandler.HardDeleteUserHandler, middlewares.CookieAuthMiddleware)
}
