package application

import (
	"GoCommerceX/auth-service/internal/domain"
	"GoCommerceX/auth-service/internal/ports"
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrRegisterEmailRequired    = errors.New("email is required")
	ErrRegisterPasswordRequired = errors.New("password is required")
	ErrUserAlreadyExists        = errors.New("user already exists")
)

type RegisterUserInput struct {
	ID       string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	UserID string
	Email  string
}

type RegisterUserUseCase struct {
	userRepository ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewRegisterUserUseCase(userRepository ports.UserRepository, passwordHasher ports.PasswordHasher) RegisterUserUseCase {
	return RegisterUserUseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (uc RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	if email == "" {
		return RegisterUserOutput{}, ErrRegisterEmailRequired
	}

	if password == "" {
		return RegisterUserOutput{}, ErrRegisterPasswordRequired
	}

	existingUser, err := uc.userRepository.FindByEmail(ctx, email)
	if err == nil && existingUser.ID != "" {
		return RegisterUserOutput{}, ErrUserAlreadyExists
	}
	passwordHash, err := uc.passwordHasher.Hash(ctx, password)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	now := time.Now().UTC()

	user := domain.User{
		ID:           strings.TrimSpace(input.ID),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := user.Validate(); err != nil {
		return RegisterUserOutput{}, err
	}

	if err := uc.userRepository.Save(ctx, user); err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
