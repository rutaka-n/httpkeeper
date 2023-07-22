package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(secret []byte, service, client string, validUntil time.Time) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(validUntil),
		Issuer:    service,
		Audience:  []string{client},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

type Payload struct {
	ExpiresAt time.Time
	Issuer    string
	Audience  []string
}

func Validate(secret []byte, tokenString string) (Payload, error) {
	payload := Payload{}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		payload.ExpiresAt = claims.ExpiresAt.Time
		payload.Issuer = claims.Issuer
		payload.Audience = claims.Audience
		return payload, nil
	} else {
		return payload, fmt.Errorf("couldn't parse calims in token %s", tokenString)
	}
}
