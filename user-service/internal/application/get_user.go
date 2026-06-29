package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/user-service/internal/domain"
	"GoCommerceX/user-service/internal/ports"
)

var (
	ErrGetUserIDRequired    = errors.New("get user: id is required")
	ErrGetUserEmailRequired = errors.New("get user by email: email is required")
	ErrUserNotFound         = errors.New("get user: user not found")
)

type GetUserInput struct {
	ID string
}

type GetUserOutput struct {
	User domain.User
}

type GetUserUseCase struct {
	userRepo ports.UserRepository
}

func NewGetUserUseCase(userRepo ports.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: userRepo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (GetUserOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return GetUserOutput{}, ErrGetUserIDRequired
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return GetUserOutput{}, ErrUserNotFound
		}
		return GetUserOutput{}, err
	}

	return GetUserOutput{User: user}, nil
}

type GetUserByEmailInput struct {
	Email string
}

type GetUserByEmailOutput struct {
	User domain.User
}

type GetUserByEmailUseCase struct {
	userRepo ports.UserRepository
}

func NewGetUserByEmailUseCase(userRepo ports.UserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{userRepo: userRepo}
}

func (uc *GetUserByEmailUseCase) Execute(ctx context.Context, input GetUserByEmailInput) (GetUserByEmailOutput, error) {
	email := strings.TrimSpace(input.Email)
	if email == "" {
		return GetUserByEmailOutput{}, ErrGetUserEmailRequired
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return GetUserByEmailOutput{}, ErrUserNotFound
		}
		return GetUserByEmailOutput{}, err
	}

	return GetUserByEmailOutput{User: user}, nil
}
