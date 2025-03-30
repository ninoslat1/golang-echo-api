package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/argon2"
)

func genSecureCookies() (string, string) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err.Error()
	}

	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return "", err.Error()
	}

	deriveKeys := argon2.IDKey(
		tokenBytes,
		salt,
		3,
		64*1024,
		4,
		32,
	)

	return base64.URLEncoding.EncodeToString(deriveKeys), ""
}

func SetSecureCookies(c echo.Context) (*http.Cookie, error) {
	token, err := genSecureCookies()
	if err != "" {
		return nil, errors.New(err)
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(4 * 7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// Set the cookie
	c.SetCookie(cookie)

	// Return the underlying http.Cookie
	return cookie, nil
}
