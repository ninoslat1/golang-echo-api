package utils

import (
	models "echo-api/models"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func genSecureCookies(userName string) (string, error) {
	cfg := &models.JwtConfig{
		Secret: os.Getenv("SECRET_KEY"),
	}

	// fmt.Println("genSecureCookies - SECRET_KEY:", cfg.Secret)

	claims := jwt.MapClaims{
		"user_name": userName,
		"exp":       time.Now().Add(4 * 7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func SetSecureCookies(c echo.Context, loginReq *models.LoginRequest) (*http.Cookie, error) {
	token, err := genSecureCookies(loginReq.UserName)
	if err != nil {
		log.Errorf("Failed to generate token for user %s: %v", loginReq.UserName, err)
		return nil, err
	}

	// Create and set the cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(4 * 7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	c.SetCookie(cookie)
	return cookie, nil
}

func ValidateSessionToken(cookie *http.Cookie) (string, error) {
	cfg := &models.JwtConfig{
		Secret: os.Getenv("SECRET_KEY"),
	}

	// fmt.Println("SECRET_KEY:", cfg.Secret)

	// Parse the token
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid session token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Extract user_code
	userCode, ok := claims["user_name"].(string)
	if !ok {
		return "", errors.New("user code missing in token")
	}

	return userCode, nil
}
