package units

import (
	"echo-api/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBase64Password(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		expectErr error
	}{
		{
			name:      "Valid Input",
			text:      "password123",
			expectErr: nil,
		},
		{
			name:      "Empty Input",
			text:      "",
			expectErr: errors.New("No encoded text"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.EncodeToBase64Password(tt.text)

			if tt.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectErr.Error())
			}
		})
	}
}
