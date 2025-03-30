package handlers

import (
	authmodels "echo-api/models"
	"echo-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService authmodels.AuthService
}

func NewAuthHandler(authService authmodels.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) LoginHandler(c echo.Context) error {
	log := logrus.New()

	loginReq := new(authmodels.LoginRequest)
	if err := c.Bind(loginReq); err != nil {
		log.Error("Error binding request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": "Invalid request"})
	}

	loginResponse, err := h.authService.Login("bromousr", loginReq)
	if err != nil {
		log.Info("Login failed: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"ERROR": err.Error()})
	}

	// Generate Secure Cookie
	cookie, err := utils.SetSecureCookies(c)
	if err != nil {
		log.Error("Failed to create session cookie for user:", loginResponse.Message)
		return c.JSON(http.StatusInternalServerError, map[string]string{"ERROR": "Failed to generate session"})
	}

	loginResponse.Cookie = cookie.Value

	return c.JSON(http.StatusOK, loginResponse)
}

func (h *AuthHandler) VerifyEmailHandler(c echo.Context) error {

	var req authmodels.VerifyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": err.Error()})
	}

	_, err := h.authService.VerifyUser("bromousr", req.Email, req.SecurityCode)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"ERROR": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Verification successful! You can now log in."})
}

func (h *AuthHandler) RegisterUserHandler(c echo.Context) error {
	var req authmodels.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": err.Error()})
	}

	if req.UserName == "" || req.Password == "" || req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": "Username, password, and email are required"})
	}

	response, err := h.authService.Register("bromousr", &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"ERROR": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}
