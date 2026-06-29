package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserIDRequired       = errors.New("user: id is required")
	ErrUserEmailRequired    = errors.New("user: email is required")
	ErrUserFirsNameRequired = errors.New("user: first name is required")
	ErrUserLastNameRequired = errors.New("user: last name is required")
)

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, firstName, lastName, phone string) User {
	now := time.Now()
	return User{
		ID:        uuid.NewString(),
		Email:     strings.TrimSpace(email),
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Phone:     strings.TrimSpace(phone),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
func (u User) Validate() error {
	if u.ID == "" {
		return ErrUserIDRequired
	}
	if u.Email == "" {
		return ErrUserEmailRequired
	}
	if u.FirstName == "" {
		return ErrUserFirsNameRequired
	}
	if u.LastName == "" {
		return ErrUserLastNameRequired
	}
	return nil
}
