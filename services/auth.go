package services

import (
	models "echo-api/models"
	repo "echo-api/repositories"
	"echo-api/utils"
	"errors"
	"fmt"
)

type authService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) models.AuthService {
	return &authService{userRepo}
}

func (s *authService) Login(dbName string, loginReq *models.LoginRequest) (*models.LoginResponse, error) {
	if loginReq.UserName == "" || loginReq.Password == "" {
		return nil, errors.New("Username and password required")
	}

	encodePassword := utils.EncodeToBase64Password(loginReq.Password)

	user, err := s.userRepo.FindByUsernameAndPassword(dbName, loginReq.UserName, encodePassword)
	if err != nil {
		return nil, errors.New("Invalid credentials")
	}

	if user.LogIn == 0 {
		return &models.LoginResponse{
			Message: fmt.Sprintf("Account %s not verified, please verify your account first", user.UserCode),
		}, nil
	}

	return &models.LoginResponse{
		Message: "Welcome " + user.UserCode,
	}, nil
}

func (s *authService) Register(dbName string, registerReq *models.RegisterRequest) (*models.RegisterResponse, error) {
	if registerReq.UserName == "" || registerReq.Password == "" || registerReq.UserCode == "" {
		return nil, errors.New("Username and password required")
	}

	encodePassword := utils.EncodeToBase64Password(registerReq.Password)

	registerReq.Password = encodePassword
	success, err := s.userRepo.RegisterUser(dbName, registerReq)
	if err != nil || !success {
		return nil, errors.New(err.Error())
	}

	response := &models.RegisterResponse{
		Message: "User registered successfully. Please verify your email.",
	}
	return response, nil
}

func (s *authService) VerifyUser(dbName, email, securityCode string) (bool, error) {
	if email == "" || securityCode == "" {
		return false, errors.New("Email and security code are required")
	}

	verified, err := s.userRepo.VerifyUser(dbName, email, securityCode)
	if err != nil {
		return false, err
	}

	if !verified {
		return false, errors.New("Invalid verification code or user already verified")
	}

	return true, nil
}
