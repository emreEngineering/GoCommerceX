package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserIDRequired       = errors.New("user: id is required")
	ErrUserEmailRequired    = errors.New("user: email is required")
	ErrPasswordHashRequired = errors.New("user: password hash is required")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, passwordHash string) User {
	now := time.Now()
	return User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u User) Validate() error {
	if u.ID == "" {
		return ErrUserIDRequired
	}
	if u.Email == "" {
		return ErrUserEmailRequired
	}
	if u.PasswordHash == "" {
		return ErrPasswordHashRequired
	}
	return nil
}
