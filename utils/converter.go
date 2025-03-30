package utils

import (
	"encoding/base64"
	"strings"
)

func EncodeToBase64Password(text string) string {
	specialChar := strings.Join(strings.Split(text, ""), "\u0000")
	return base64.StdEncoding.EncodeToString([]byte(specialChar))
}

func DecodePasswordFromBase64(text string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(decode), "\u0000", ""), nil
}
