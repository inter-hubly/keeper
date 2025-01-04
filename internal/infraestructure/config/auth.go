package config

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-strong-secret-key")

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Tenant   uint64 `json:"tenant"`
	jwt.RegisteredClaims
}

func GenerateBearerToken(_ context.Context, username string, clientId uint64) (string, error) {
	claims := CustomClaims{
		Username: username,
		Tenant:   clientId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
