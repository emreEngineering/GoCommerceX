package application

import (
	"context"
	"errors"
	"testing"

	"GoCommerceX/auth-service/internal/domain"
)

type fakeUserRepository struct {
	usersByEmail map[string]domain.User
	savedUser    domain.User
	saveErr      error
	findErr      error
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByEmail: make(map[string]domain.User),
	}
}

func (r *fakeUserRepository) Save(ctx context.Context, user domain.User) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	r.savedUser = user
	r.usersByEmail[user.Email] = user

	return nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if r.findErr != nil {
		return domain.User{}, r.findErr
	}

	user, ok := r.usersByEmail[email]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil
}

type fakePasswordHasher struct {
	hashValue string
	hashErr   error
}

func (h fakePasswordHasher) Hash(ctx context.Context, plainPassword string) (string, error) {
	if h.hashErr != nil {
		return "", h.hashErr
	}

	return h.hashValue, nil
}

func (h fakePasswordHasher) Compare(ctx context.Context, plainPassword string, passwordHash string) error {
	return nil
}

func TestRegisterUserUseCaseExecute(t *testing.T) {
	t.Run("registers user successfully", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{
			hashValue: "hashed-password",
		}

		useCase := NewRegisterUserUseCase(repository, hasher)

		output, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:       "user-1",
			Email:    " user@example.com ",
			Password: " secret-password ",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if output.UserID != "user-1" {
			t.Fatalf("expected user id user-1, got %s", output.UserID)
		}

		if output.Email != "user@example.com" {
			t.Fatalf("expected trimmed email, got %s", output.Email)
		}

		if repository.savedUser.PasswordHash != "hashed-password" {
			t.Fatalf("expected saved password hash, got %s", repository.savedUser.PasswordHash)
		}
	})

	t.Run("returns error when email is missing", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{hashValue: "hashed-password"}

		useCase := NewRegisterUserUseCase(repository, hasher)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:       "user-1",
			Password: "secret-password",
		})

		if !errors.Is(err, ErrRegisterEmailRequired) {
			t.Fatalf("expected error %v, got %v", ErrRegisterEmailRequired, err)
		}
	})

	t.Run("returns error when password is missing", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{hashValue: "hashed-password"}

		useCase := NewRegisterUserUseCase(repository, hasher)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:    "user-1",
			Email: "user@example.com",
		})

		if !errors.Is(err, ErrRegisterPasswordRequired) {
			t.Fatalf("expected error %v, got %v", ErrRegisterPasswordRequired, err)
		}
	})

	t.Run("returns error when user already exists", func(t *testing.T) {
		repository := newFakeUserRepository()
		repository.usersByEmail["user@example.com"] = domain.User{
			ID:           "existing-user",
			Email:        "user@example.com",
			PasswordHash: "existing-hash",
		}

		hasher := fakePasswordHasher{hashValue: "hashed-password"}

		useCase := NewRegisterUserUseCase(repository, hasher)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:       "user-1",
			Email:    "user@example.com",
			Password: "secret-password",
		})

		if !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected error %v, got %v", ErrUserAlreadyExists, err)
		}
	})

	t.Run("returns error when password hashing fails", func(t *testing.T) {
		repository := newFakeUserRepository()
		hashErr := errors.New("hash failed")

		hasher := fakePasswordHasher{
			hashErr: hashErr,
		}

		useCase := NewRegisterUserUseCase(repository, hasher)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:       "user-1",
			Email:    "user@example.com",
			Password: "secret-password",
		})

		if !errors.Is(err, hashErr) {
			t.Fatalf("expected error %v, got %v", hashErr, err)
		}
	})

	t.Run("returns error when save fails", func(t *testing.T) {
		saveErr := errors.New("save failed")

		repository := newFakeUserRepository()
		repository.saveErr = saveErr

		hasher := fakePasswordHasher{
			hashValue: "hashed-password",
		}

		useCase := NewRegisterUserUseCase(repository, hasher)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:       "user-1",
			Email:    "user@example.com",
			Password: "secret-password",
		})

		if !errors.Is(err, saveErr) {
			t.Fatalf("expected error %v, got %v", saveErr, err)
		}
	})
}
