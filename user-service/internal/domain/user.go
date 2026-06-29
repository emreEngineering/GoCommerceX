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

func NewUser(id, email, firstName, lastName, phone string) User {
	now := time.Now()
	if strings.TrimSpace(id) == "" {
		id = uuid.NewString()
	}
	return User{
		ID:        strings.TrimSpace(id),
		Email:     strings.TrimSpace(email),
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Phone:     strings.TrimSpace(phone),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return ErrUserIDRequired
	}
	if strings.TrimSpace(u.Email) == "" {
		return ErrUserEmailRequired
	}
	if strings.TrimSpace(u.FirstName) == "" {
		return ErrUserFirsNameRequired
	}
	if strings.TrimSpace(u.LastName) == "" {
		return ErrUserLastNameRequired
	}
	return nil
}
