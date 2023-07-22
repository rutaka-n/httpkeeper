package token

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecret(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
