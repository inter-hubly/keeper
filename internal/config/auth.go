package config

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/server"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	LoggedUser hctx.Logged `json:"loggedUser"`
	jwt.RegisteredClaims
}

func GenerateBearerToken(_ context.Context, username, userId, tenantId string) (string, error) {
	claims := CustomClaims{
		LoggedUser: hctx.Logged{
			Username: username,
			Tenant:   tenantId,
			UserId:   userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(server.GetEnvironment().HashEncrypt))
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckHashPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
