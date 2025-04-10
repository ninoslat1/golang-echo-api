package handlers

import (
	"echo-api/views"

	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	return views.Home().Render(c.Request().Context(), c.Response().Writer)
}
