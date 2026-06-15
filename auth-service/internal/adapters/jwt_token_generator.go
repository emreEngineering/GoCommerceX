package adapters

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	secretKey []byte
}

func NewJWTTokenGenerator(secretKey string) JWTTokenGenerator {
	return JWTTokenGenerator{secretKey: []byte(secretKey)}
}

func (g JWTTokenGenerator) Generate(ctx context.Context, userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.secretKey)
}
