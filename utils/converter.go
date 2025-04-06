package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

func EncodeToBase64Password(text string) (string, error) {
	if text == "" {
		return "", errors.New("No encoded text")
	}

	specialChar := strings.Join(strings.Split(text, ""), "\u0000")
	return base64.StdEncoding.EncodeToString([]byte(specialChar)), nil
}

func DecodePasswordFromBase64(text string) (string, error) {
	if text == "" {
		return "", errors.New("No decoded text")
	}

	decode, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", errors.New("Encoded text not valid")
	}

	return strings.ReplaceAll(string(decode), "\u0000", ""), nil
}
