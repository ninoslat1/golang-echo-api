package handlers

import (
	authmodels "echo-api/models"
	"echo-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService authmodels.AuthService
	log         *logrus.Logger
}

func NewAuthHandler(authService authmodels.AuthService, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{authService, log}
}

func (h *AuthHandler) LoginHandler(c echo.Context) error {
	loginReq := new(authmodels.LoginRequest)
	if err := c.Bind(loginReq); err != nil {
		log.Error("Error binding request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": "Invalid request"})
	}

	user, err := h.authService.Login("bromousr", loginReq)
	if err != nil {
		log.Info("Login failed: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"ERROR": err.Error()})
	}

	cookie, err := utils.SetSecureCookies(c, loginReq)
	if err != nil {
		log.Error("Failed to create session cookie for user:", user.UserName)
		return c.JSON(http.StatusInternalServerError, map[string]string{"ERROR": "Failed to generate session"})
	}

	user.UserName = cookie.Value

	response := authmodels.LoginResponse{
		Message: "Welcome " + loginReq.UserName,
		Cookie:  cookie.Value,
	}

	return c.JSON(http.StatusOK, response)
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

func (h *AuthHandler) SoftDeleteUserHandler(c echo.Context) error {
	loginReq := new(authmodels.LoginRequest)
	if err := c.Bind(loginReq); err != nil {
		log.Error("Error binding request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": "Invalid request"})
	}

	err := h.authService.SoftDeleteUser("bromousr", loginReq)
	if err != nil {
		log.Error("Soft delete failed: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"ERROR": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User soft deleted successfully"})
}

func (h *AuthHandler) HardDeleteUserHandler(c echo.Context) error {
	userName, ok := c.Get("user_name").(string)
	if !ok || userName == "" {
		log.Error("Unauthorized access: Missing user_name from context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"ERROR": "Unauthorized"})
	}

	loginReq := new(authmodels.LoginRequest)
	if err := c.Bind(loginReq); err != nil {
		log.Error("Error binding request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"ERROR": "Invalid request"})
	}

	if userName != loginReq.UserName {
		log.Error("Unauthorized attempt to delete another user's account")
		return c.JSON(http.StatusForbidden, map[string]string{"ERROR": "You can only delete your own account"})
	}

	err := h.authService.HardDeleteUser("bromousr", loginReq)
	if err != nil {
		log.Error("Soft delete failed: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"ERROR": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User permanent delete successfully"})
}
