package main

import (
	router "echo-api/router"
	service "echo-api/services"

	"github.com/labstack/echo/v4"
)

func main() {
	log := service.InitLogger()
	e := echo.New()

	router.SetupRoutes(e, log)

	e.Logger.Fatal(e.Start(":1323"))
}
