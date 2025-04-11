package handlers

import (
	"echo-api/views"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return views.LoginPage().Render(c.Request().Context(), c.Response().Writer)
}
