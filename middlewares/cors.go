package middlewares

import (
	"echo-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
}

func CookieAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Retrieve cookie
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
