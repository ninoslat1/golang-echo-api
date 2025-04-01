package main

import (
	middleware "echo-api/middlewares"
	router "echo-api/router"
	service "echo-api/services"

	"github.com/labstack/echo/v4"
)

func main() {
	log := service.InitLogger()

	e := echo.New()
	middleware.SetupMiddlewares(e)

	router.SetupRoutes(e, log)

	e.Logger.Fatal(e.Start(":1323"))
}
