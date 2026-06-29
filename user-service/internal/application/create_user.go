package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/user-service/internal/domain"
	"GoCommerceX/user-service/internal/ports"
)

// Hatalar
var (
	ErrCreateUserEmailRequired     = errors.New("create user: email is required")
	ErrCreateUserFirstNameRequired = errors.New("create user: first name is required")
	ErrCreateUserLastNameRequired  = errors.New("create user: last name is required")
	ErrUserAlreadyExists           = errors.New("create user: user already exists")
)

// Girdi
type CreateUserInput struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

// Çıktı
type CreateUserOutput struct {
	User domain.User
}

// Use Case
type CreateUserUseCase struct {
	userRepo ports.UserRepository
}

func NewCreateUserUseCase(userRepo ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	// Temizle
	email := strings.TrimSpace(input.Email)
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)
	phone := strings.TrimSpace(input.Phone)

	// Kontroller
	if email == "" {
		return CreateUserOutput{}, ErrCreateUserEmailRequired
	}
	if firstName == "" {
		return CreateUserOutput{}, ErrCreateUserFirstNameRequired
	}
	if lastName == "" {
		return CreateUserOutput{}, ErrCreateUserLastNameRequired
	}

	// Email zaten var mı?
	_, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return CreateUserOutput{}, ErrUserAlreadyExists
	}
	if !errors.Is(err, ports.ErrUserNotFound) {
		return CreateUserOutput{}, err
	}

	// Domain objesi oluştur
	user := domain.NewUser(input.ID, email, firstName, lastName, phone)
	if err := user.Validate(); err != nil {
		return CreateUserOutput{}, err
	}

	// Kaydet
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{User: user}, nil
}
