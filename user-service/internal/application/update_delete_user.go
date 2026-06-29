package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/user-service/internal/domain"
	"GoCommerceX/user-service/internal/ports"
)

var (
	ErrUpdateUserIDRequired        = errors.New("update user: id is required")
	ErrUpdateUserFirstNameRequired = errors.New("update user: first name is required")
	ErrUpdateUserLastNameRequired  = errors.New("update user: last name is required")
	ErrDeleteUserIDRequired        = errors.New("delete user: id is required")
)

// ========== UpdateUserUseCase ==========

type UpdateUserInput struct {
	ID        string
	FirstName string
	LastName  string
	Phone     string
}

type UpdateUserOutput struct {
	User domain.User
}

type UpdateUserUseCase struct {
	userRepo ports.UserRepository
}

func NewUpdateUserUseCase(userRepo ports.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (UpdateUserOutput, error) {
	id := strings.TrimSpace(input.ID)
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)
	phone := strings.TrimSpace(input.Phone)

	if id == "" {
		return UpdateUserOutput{}, ErrUpdateUserIDRequired
	}
	if firstName == "" {
		return UpdateUserOutput{}, ErrUpdateUserFirstNameRequired
	}
	if lastName == "" {
		return UpdateUserOutput{}, ErrUpdateUserLastNameRequired
	}

	existingUser, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return UpdateUserOutput{}, ErrUserNotFound
		}
		return UpdateUserOutput{}, err
	}

	existingUser.FirstName = firstName
	existingUser.LastName = lastName
	existingUser.Phone = phone

	if err := uc.userRepo.Update(ctx, existingUser); err != nil {
		return UpdateUserOutput{}, err
	}

	return UpdateUserOutput{User: existingUser}, nil
}

// ========== DeleteUserUseCase ==========

type DeleteUserInput struct {
	ID string
}

type DeleteUserOutput struct {
	Success bool
}

type DeleteUserUseCase struct {
	userRepo ports.UserRepository
}

func NewDeleteUserUseCase(userRepo ports.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepo: userRepo}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (DeleteUserOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return DeleteUserOutput{}, ErrDeleteUserIDRequired
	}

	if err := uc.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return DeleteUserOutput{}, ErrUserNotFound
		}
		return DeleteUserOutput{}, err
	}

	return DeleteUserOutput{Success: true}, nil
}
