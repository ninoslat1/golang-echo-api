package models

type ValidatorTestCase[T any] struct {
	Name      string
	Input     T
	ExpectErr error
}
