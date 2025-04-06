package utils

import (
	"echo-api/models"
	"errors"
	"testing"
)

func TestLoginRequest(t *testing.T) {
	testCases := []models.ValidatorTestCase[*models.LoginRequest]{
		{
			Name:      "Valid Input",
			Input:     &models.LoginRequest{UserName: "nino", Password: "password123"},
			ExpectErr: nil,
		},
		{
			Name:      "Empty Username",
			Input:     &models.LoginRequest{UserName: "", Password: "password123"},
			ExpectErr: errors.New("Username and password required"),
		},
		{
			Name:      "Empty Password",
			Input:     &models.LoginRequest{UserName: "test123", Password: ""},
			ExpectErr: errors.New("Username and password required"),
		},
	}

	RunValidatorTest(t, testCases, LoginRequestValidator)
}

func TestRegisterRequest(t *testing.T) {
	testCases := []models.ValidatorTestCase[*models.RegisterRequest]{
		{
			Name:      "Valid Input and Valid Email",
			Input:     &models.RegisterRequest{UserName: "test", Password: "test", UserCode: "test", Email: "test@gmail.com"},
			ExpectErr: nil,
		},
		{
			Name:      "Null Username and Valid Email",
			Input:     &models.RegisterRequest{UserName: "", Password: "test", UserCode: "test", Email: "test@gmail.com"},
			ExpectErr: errors.New("Username, password, email, and user code are required"),
		},
		{
			Name:      "Invalid Email",
			Input:     &models.RegisterRequest{UserName: "test", Password: "test", UserCode: "test", Email: "john@.site.com"},
			ExpectErr: errors.New("Email is not valid"),
		},
		{
			Name:      "Null Email",
			Input:     &models.RegisterRequest{UserName: "test", Password: "test", UserCode: "test", Email: ""},
			ExpectErr: errors.New("Username, password, email, and user code are required"),
		},
		{
			Name:      "Null UserCode and Valid Email",
			Input:     &models.RegisterRequest{UserName: "test", Password: "test", UserCode: "", Email: "test@gmail.com"},
			ExpectErr: errors.New("Username, password, email, and user code are required"),
		},
		{
			Name:      "Null Password and Valid Email",
			Input:     &models.RegisterRequest{UserName: "test", Password: "", UserCode: "test", Email: "test@gmail.com"},
			ExpectErr: errors.New("Username, password, email, and user code are required"),
		},
	}

	RunValidatorTest(t, testCases, RegisterRequestValidator)
}

func TestVerifyUserValidator(t *testing.T) {
	testCases := []models.ValidatorTestCase[*models.VerifyUserInput]{
		{
			Name:      "Invalid Input",
			Input:     &models.VerifyUserInput{Email: "", SecurityCode: ""},
			ExpectErr: errors.New("Email and security code are required"),
		},
		{
			Name:      "Null Email",
			Input:     &models.VerifyUserInput{Email: "", SecurityCode: "123456"},
			ExpectErr: errors.New("Email and security code are required"),
		},
		{
			Name:      "Null Security Code",
			Input:     &models.VerifyUserInput{Email: "test@gmail.com", SecurityCode: ""},
			ExpectErr: errors.New("Email and security code are required"),
		},
		{
			Name:      "Invalid Security Code",
			Input:     &models.VerifyUserInput{Email: "test@gmail.com", SecurityCode: "1234"},
			ExpectErr: errors.New("Security code is not valid"),
		},
		{
			Name:      "Invalid Email",
			Input:     &models.VerifyUserInput{Email: "john@.site.com", SecurityCode: "123456"},
			ExpectErr: errors.New("Email is not valid"),
		},
		{
			Name:      "Valid Input",
			Input:     &models.VerifyUserInput{Email: "test@gmail.com", SecurityCode: "123456"},
			ExpectErr: nil,
		},
	}

	RunValidatorTest(t, testCases, VerifyUserValidator)
}

func TestValidEmail(t *testing.T) {

	testCases := []models.ValidatorTestCase[string]{
		{
			Name:      "Valid Email",
			Input:     "test@example.com",
			ExpectErr: nil,
		},
		{
			Name:      "Missing @",
			Input:     "testexample.com",
			ExpectErr: errors.New("Email is not valid"),
		},
		{
			Name:      "Missing domain",
			Input:     "test@",
			ExpectErr: errors.New("Email is not valid"),
		},
		{
			Name:      "Multiple @",
			Input:     "a@b@c.com",
			ExpectErr: errors.New("Email is not valid"),
		},
	}

	RunValidatorTest(t, testCases, func(input string) error {
		if !IsValidEmail(input) {
			return errors.New("Email is not valid")
		}
		return nil
	})
}
