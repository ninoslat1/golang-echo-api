package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func GenerateSecurityCode() (string, error) {
	source := rand.NewSource(time.Now().UnixNano())
	if source == nil {
		return "", errors.New("failed to initialize random source")
	}

	rng := rand.New(source)
	if rng == nil {
		return "", errors.New("failed to create random number generator")
	}

	securityCode := fmt.Sprintf("%06d", rng.Intn(1000000))
	if len(securityCode) != 6 {
		return "", errors.New("invalid security code generated")
	}

	return securityCode, nil
}
