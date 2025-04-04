package services

import (
	models "echo-api/models"
	"echo-api/utils"
	"fmt"
)

type authService struct {
	userRepo models.UserRepository
}

func NewAuthService(userRepo models.UserRepository) models.AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(dbName string, registerReq *models.RegisterRequest) (*models.RegisterResponse, error) {
	// Validasi data masuk
	if err := utils.RegisterRequestValidator(registerReq); err != nil {
		return nil, err
	}

	// Encode password sebelum disimpan
	registerReq.Password, _ = utils.EncodeToBase64Password(registerReq.Password)

	// Generate kode keamanan
	securityCode, err := utils.GenerateSecurityCode()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate security code: %w", err)
	}

	registerReq.SecurityCode = securityCode
	registerReq.LogIn = 0

	err = s.userRepo.RegisterUser(dbName, registerReq)
	if err != nil {
		return nil, err
	}

	if err := utils.SendVerificationEmail(registerReq.Email, securityCode); err != nil {
		return nil, fmt.Errorf("Failed to send verification email: %w", err)
	}

	return &models.RegisterResponse{Message: "User registered successfully. Please verify your email."}, nil
}

func (s *authService) VerifyUser(dbName, email, securityCode string) (bool, error) {
	if err := utils.VerifyUserValidator(email, securityCode); err != nil {
		return false, err
	}

	return s.userRepo.VerifyUser(dbName, email, securityCode)
}

func (s *authService) Login(dbName string, loginReq *models.LoginRequest) (*models.User, error) {
	if err := utils.LoginRequestValidator(loginReq); err != nil {
		return nil, err
	}

	encodedPassword, _ := utils.EncodeToBase64Password(loginReq.Password)

	user, err := s.userRepo.FindByUsernameAndPassword(dbName, loginReq.UserName, encodedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) ResendVerifyCode(dbName, email string) error {
	securityCode, err := utils.GenerateSecurityCode()
	if err != nil {
		return fmt.Errorf("Failed to generate security code: %w", err)
	}

	// Perbarui kode verifikasi di database
	err = s.userRepo.ResendVerifyCode(dbName, email, securityCode)
	if err != nil {
		return err
	}

	// Kirim ulang email verifikasi
	err = utils.SendVerificationEmail(email, securityCode)
	if err != nil {
		return fmt.Errorf("Failed to send verification email: %w", err)
	}

	return nil
}

func (s *authService) SoftDeleteUser(dbName string, loginReq *models.LoginRequest) error {
	if err := utils.LoginRequestValidator(loginReq); err != nil {
		return err
	}

	encodedPassword, err := utils.EncodeToBase64Password(loginReq.Password)

	if err != nil {
		return err
	}

	err = s.userRepo.SoftDeleteUser(dbName, loginReq.UserName, encodedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) HardDeleteUser(dbName string, loginReq *models.LoginRequest) error {
	if err := utils.LoginRequestValidator(loginReq); err != nil {
		return err
	}

	encodedPassword, err := utils.EncodeToBase64Password(loginReq.Password)

	if err != nil {
		return err
	}

	err = s.userRepo.HardDeleteUser(dbName, loginReq.UserName, encodedPassword)
	if err != nil {
		return err
	}

	return nil
}
