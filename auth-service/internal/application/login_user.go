package application

import (
	"GoCommerceX/auth-service/internal/ports"
	"context"
	"errors"
	"strings"
)

var (
	ErrLoginEmailRequired    = errors.New("email is required")
	ErrLoginPasswordRequired = errors.New("password is required")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken string
}

type LoginUserUseCase struct {
	userRepository ports.UserRepository
	passwordHasher ports.PasswordHasher
	tokenGenerator ports.TokenGenerator
}

func NewLoginUserUseCase(userRepository ports.UserRepository, passwordHasher ports.PasswordHasher, tokenGenerator ports.TokenGenerator) LoginUserUseCase {
	return LoginUserUseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (uc LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	if email == "" {
		return LoginUserOutput{}, ErrLoginEmailRequired
	}

	if password == "" {
		return LoginUserOutput{}, ErrLoginPasswordRequired
	}

	user, err := uc.userRepository.FindByEmail(ctx, email)
	if err != nil || user.ID == "" {
		return LoginUserOutput{}, ErrInvalidCredentials
	}
	if err := uc.passwordHasher.Compare(ctx, password, user.PasswordHash); err != nil {
		return LoginUserOutput{}, ErrInvalidCredentials
	}
	accessToken, err := uc.tokenGenerator.Generate(ctx, user.ID, user.Email)
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		AccessToken: accessToken,
	}, nil
}
