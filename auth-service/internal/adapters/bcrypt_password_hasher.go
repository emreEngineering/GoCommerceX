package adapters

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() BcryptPasswordHasher {
	return BcryptPasswordHasher{}
}

func (h BcryptPasswordHasher) Hash(ctx context.Context, plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h BcryptPasswordHasher) Compare(ctx context.Context, plainPassword string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
}
