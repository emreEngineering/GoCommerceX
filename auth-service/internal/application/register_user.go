package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/auth-service/internal/domain"
	"GoCommerceX/auth-service/internal/ports"
)

var (
	ErrRegisterEmailRequired    = errors.New("register: email is required")
	ErrRegisterPasswordRequired = errors.New("register: password is required")
	ErrUserAlreadyExists        = errors.New("register: user already exists")
)

type RegisterUserUseCase struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewRegisterUserUseCase(userRepo ports.UserRepository, passwordHasher ports.PasswordHasher) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

type RegisterUserInput struct {
	Email    string
	Password string
}

type RegisterUserOutput struct {
	UserID string
	Email  string
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	if email == "" {
		return RegisterUserOutput{}, ErrRegisterEmailRequired
	}
	if password == "" {
		return RegisterUserOutput{}, ErrRegisterPasswordRequired
	}

	// check existing user
	_, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return RegisterUserOutput{}, ErrUserAlreadyExists
	}
	if !errors.Is(err, ports.ErrUserNotFound) {
		return RegisterUserOutput{}, err
	}

	hashedPassword, err := uc.passwordHasher.Hash(ctx, password)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	user := domain.NewUser(email, hashedPassword)
	if err := user.Validate(); err != nil {
		return RegisterUserOutput{}, err
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
