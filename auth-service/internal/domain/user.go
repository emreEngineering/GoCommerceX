package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrUserIDRequired       = errors.New("user id is required")
	ErrUserEmailRequired    = errors.New("user email is required")
	ErrPasswordHashRequired = errors.New("password hash is required")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return ErrUserIDRequired
	}

	if strings.TrimSpace(u.Email) == "" {
		return ErrUserEmailRequired
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return ErrPasswordHashRequired
	}

	return nil
}
