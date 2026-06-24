package adapters

import (
	"context"
	"fmt"
	"time"

	"GoCommerceX/auth-service/internal/ports"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	secret string
}

func NewJWTTokenGenerator(secret string) *JWTTokenGenerator {
	return &JWTTokenGenerator{secret: secret}
}

func (g *JWTTokenGenerator) Generate(ctx context.Context, userID string, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

var _ ports.TokenGenerator = (*JWTTokenGenerator)(nil)
