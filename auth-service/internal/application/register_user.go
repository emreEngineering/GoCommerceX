package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/auth-service/internal/domain"
	"GoCommerceX/auth-service/internal/ports"
)

var (
	ErrRegisterEmailRequired     = errors.New("register: email is required")
	ErrRegisterPasswordRequired  = errors.New("register: password is required")
	ErrRegisterFirstNameRequired = errors.New("register: first name is required")
	ErrRegisterLastNameRequired  = errors.New("register: last name is required")
	ErrUserAlreadyExists         = errors.New("register: user already exists")
	ErrUserProfileAlreadyExists  = errors.New("register: user profile already exists")
	ErrUserProfileCreationFailed = errors.New("register: user profile creation failed")
)

type RegisterUserUseCase struct {
	userRepo           ports.UserRepository
	passwordHasher     ports.PasswordHasher
	userProfileCreator ports.UserProfileCreator
}

func NewRegisterUserUseCase(userRepo ports.UserRepository, passwordHasher ports.PasswordHasher, userProfileCreator ports.UserProfileCreator) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:           userRepo,
		passwordHasher:     passwordHasher,
		userProfileCreator: userProfileCreator,
	}
}

type RegisterUserInput struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}

type RegisterUserOutput struct {
	UserID string
	Email  string
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)
	phone := strings.TrimSpace(input.Phone)

	if email == "" {
		return RegisterUserOutput{}, ErrRegisterEmailRequired
	}
	if password == "" {
		return RegisterUserOutput{}, ErrRegisterPasswordRequired
	}
	if firstName == "" {
		return RegisterUserOutput{}, ErrRegisterFirstNameRequired
	}
	if lastName == "" {
		return RegisterUserOutput{}, ErrRegisterLastNameRequired
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
	if trimmedID := strings.TrimSpace(input.ID); trimmedID != "" {
		user.ID = trimmedID
	}
	if err := user.Validate(); err != nil {
		return RegisterUserOutput{}, err
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return RegisterUserOutput{}, err
	}

	if err := uc.userProfileCreator.Create(ctx, ports.UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
	}); err != nil {
		_ = uc.userRepo.Delete(ctx, user.ID)
		if errors.Is(err, ErrUserProfileAlreadyExists) {
			return RegisterUserOutput{}, ErrUserProfileAlreadyExists
		}
		return RegisterUserOutput{}, ErrUserProfileCreationFailed
	}

	return RegisterUserOutput{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
