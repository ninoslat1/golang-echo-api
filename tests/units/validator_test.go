package units

import (
	"echo-api/models"
	"echo-api/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginRequest(t *testing.T) {
	tests := []struct {
		name      string
		loginReq  *models.LoginRequest
		expectErr error
	}{
		{
			name:      "Valid Input",
			loginReq:  &models.LoginRequest{UserName: "nino", Password: "password123"},
			expectErr: nil,
		},
		{
			name:      "Empty Username",
			loginReq:  &models.LoginRequest{UserName: "", Password: "password123"},
			expectErr: errors.New("Username and password required"),
		},
		{
			name:      "Empty Password",
			loginReq:  &models.LoginRequest{UserName: "nino", Password: ""},
			expectErr: errors.New("Username and password required"),
		},
		{
			name:      "Both Empty",
			loginReq:  &models.LoginRequest{UserName: "", Password: ""},
			expectErr: errors.New("Username and password required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.LoginRequestValidator(tt.loginReq)

			if tt.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectErr.Error())
			}
		})
	}
}
