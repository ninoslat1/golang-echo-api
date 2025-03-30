package utils

import (
	"echo-api/models"
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func LoginRequestValidator(loginReq *models.LoginRequest) error {
	if loginReq.UserName == "" || loginReq.Password == "" {
		return errors.New("Username and password required")
	}
	return nil
}

func RegisterRequestValidator(registerReq *models.RegisterRequest) error {
	if registerReq.UserName == "" || registerReq.Password == "" || registerReq.UserCode == "" || registerReq.Email == "" {
		return errors.New("Username, password, email, and user code are required")
	}

	return nil
}

func VerifyUserValidator(email, securityCode string) error {
	if email == "" || securityCode == "" {
		return errors.New("Email and security code are required")
	}

	return nil
}
