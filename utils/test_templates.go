package utils

import (
	"echo-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunValidatorTest[T any](t *testing.T, cases []models.ValidatorTestCase[T], validatorFunc func(T) error) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validatorFunc(tc.Input)

			if tc.ExpectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.ExpectErr.Error())
			}
		})
	}
}
