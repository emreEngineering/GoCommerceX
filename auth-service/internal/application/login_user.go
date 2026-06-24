package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/auth-service/internal/ports"
)

var (
	ErrLoginEmailRequired    = errors.New("login: email is required")
	ErrLoginPasswordRequired = errors.New("login: password is required")
	ErrInvalidCredentials    = errors.New("login: invalid credentials")
)

type LoginUserUseCase struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
	tokenGenerator ports.TokenGenerator
}

func NewLoginUserUseCase(userRepo ports.UserRepository, passwordHasher ports.PasswordHasher, tokenGenerator ports.TokenGenerator) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken string
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	if email == "" {
		return LoginUserOutput{}, ErrLoginEmailRequired
	}
	if password == "" {
		return LoginUserOutput{}, ErrLoginPasswordRequired
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return LoginUserOutput{}, ErrInvalidCredentials
		}
		return LoginUserOutput{}, err
	}

	if err := uc.passwordHasher.Compare(ctx, password, user.PasswordHash); err != nil {
		return LoginUserOutput{}, ErrInvalidCredentials
	}

	token, err := uc.tokenGenerator.Generate(ctx, user.ID, user.Email)
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{AccessToken: token}, nil
}
