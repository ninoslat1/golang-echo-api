package utils

import (
	"echo-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunValidatorTest[T any](t *testing.T, cases []models.ValidatorTestCase[T], validatorFunc any) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var result string
			var err error

			switch vf := validatorFunc.(type) {
			case func(T) error:
				err = vf(tc.Input)
			case func(T) (string, error):
				result, err = vf(tc.Input)
				if tc.ExpectErr == nil {
					assert.NotEmpty(t, result)
				} else {
					assert.Empty(t, result)
				}
			default:
				t.Fatalf("Unsupported validator function signature")
			}

			if tc.ExpectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.ExpectErr.Error())
			}
		})
	}
}
