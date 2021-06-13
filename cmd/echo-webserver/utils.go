package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomString(len int) (string, error) {
	buf := make([]byte, len)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("error: could not rand.Read: %w", err)
	}
	return base64.RawStdEncoding.EncodeToString(buf), nil
}
