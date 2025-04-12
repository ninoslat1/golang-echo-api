package middlewares

import (
	"echo-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CookieAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: No session found"})
		}

		userName, err := utils.ValidateSessionToken(cookie)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: " + err.Error()})
		}

		c.Set("user_name", userName)
		return next(c)
	}
}

func CookiePageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		_, err = utils.ValidateSessionToken(cookie)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}
