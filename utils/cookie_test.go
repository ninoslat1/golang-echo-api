package utils

import (
	"database/sql"
	models "echo-api/models"
	"errors"
	"os"
	"testing"
)

func TestGenerateSecureCookies(t *testing.T) {
	os.Setenv("SECRET_KEY", "testingsecret")

	testCases := []models.ValidatorTestCase[*models.User]{
		{
			Name: "Valid user",
			Input: &models.User{
				UserName: "testuser",
				LogIn:    sql.NullInt32{Int32: 1, Valid: true},
			},
			ExpectErr: nil,
		},
		{
			Name:      "Not valid user",
			Input:     nil,
			ExpectErr: errors.New("User is not exist"),
		},
	}

	RunValidatorTest(t, testCases, genSecureCookies)
}
