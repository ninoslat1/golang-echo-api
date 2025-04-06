package utils

import (
	models "echo-api/models"
	"errors"
	"testing"
)

func TestEncodeBase64Password(t *testing.T) {
	testCases := []models.ValidatorTestCase[string]{
		{
			Name:      "Valid Input",
			Input:     "password123",
			ExpectErr: nil,
		},
		{
			Name:      "Empty Input",
			Input:     "",
			ExpectErr: errors.New("No encoded text"),
		},
	}

	RunValidatorTest(t, testCases, func(input string) error {
		_, err := EncodeToBase64Password(input)
		return err
	})
}

func TestDecodeBase64Password(t *testing.T) {
	testCases := []models.ValidatorTestCase[string]{
		{
			Name:      "Valid Base64 Input",
			Input:     "cGFzc3dvcmQxMjM=",
			ExpectErr: nil,
		},
		{
			Name:      "Empty Input",
			Input:     "",
			ExpectErr: errors.New("No decoded text"),
		},
		{
			Name:      "Invalid Base64 Input",
			Input:     "invalid@@@",
			ExpectErr: errors.New("Encoded text not valid"),
		},
	}

	RunValidatorTest(t, testCases, func(input string) error {
		_, err := DecodePasswordFromBase64(input)
		return err
	})
}
