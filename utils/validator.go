package utils

import (
	"echo-api/models"
	"errors"
	"strings"

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

	isValid := IsValidEmail(registerReq.Email)

	if !isValid {
		return errors.New("Email is not valid")
	}

	return nil
}

func VerifyUserValidator(input *models.VerifyUserInput) error {
	if input.Email == "" || input.SecurityCode == "" {
		return errors.New("Email and security code are required")
	}

	if len(input.SecurityCode) != 6 {
		return errors.New("Security code is not valid")
	}

	isValid := IsValidEmail(input.Email)

	if !isValid {
		return errors.New("Email is not valid")
	}

	return nil
}

func IsValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	user := parts[0]
	domain := parts[1]

	if len(user) == 0 {
		return false
	}

	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	if strings.Contains(domain, "..") {
		return false
	}

	for _, r := range email {
		if !(r >= 'a' && r <= 'z' ||
			r >= 'A' && r <= 'Z' ||
			r >= '0' && r <= '9' ||
			r == '@' || r == '.' || r == '-' || r == '_' || r == '+') {
			return false
		}
	}

	return true
}
