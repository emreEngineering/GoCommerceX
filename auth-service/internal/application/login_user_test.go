package application

import (
	"context"
	"errors"
	"testing"

	"GoCommerceX/auth-service/internal/domain"
)

type fakeTokenGenerator struct {
	token       string
	generateErr error
}

func (g fakeTokenGenerator) Generate(ctx context.Context, userID string, email string) (string, error) {
	if g.generateErr != nil {
		return "", g.generateErr
	}

	return g.token, nil
}

func TestLoginUserUseCaseExecute(t *testing.T) {
	t.Run("logs in user successfully", func(t *testing.T) {
		repository := newFakeUserRepository()
		repository.usersByEmail["user@example.com"] = domain.User{
			ID:           "user-1",
			Email:        "user@example.com",
			PasswordHash: "hashed-password",
		}

		hasher := fakePasswordHasher{}
		tokenGenerator := fakeTokenGenerator{
			token: "access-token",
		}

		useCase := NewLoginUserUseCase(repository, hasher, tokenGenerator)

		output, err := useCase.Execute(context.Background(), LoginUserInput{
			Email:    " user@example.com ",
			Password: " secret-password ",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if output.AccessToken != "access-token" {
			t.Fatalf("expected access token, got %s", output.AccessToken)
		}
	})

	t.Run("returns error when email is missing", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{}
		tokenGenerator := fakeTokenGenerator{token: "access-token"}

		useCase := NewLoginUserUseCase(repository, hasher, tokenGenerator)

		_, err := useCase.Execute(context.Background(), LoginUserInput{
			Password: "secret-password",
		})

		if !errors.Is(err, ErrLoginEmailRequired) {
			t.Fatalf("expected error %v, got %v", ErrLoginEmailRequired, err)
		}
	})

	t.Run("returns error when password is missing", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{}
		tokenGenerator := fakeTokenGenerator{token: "access-token"}

		useCase := NewLoginUserUseCase(repository, hasher, tokenGenerator)

		_, err := useCase.Execute(context.Background(), LoginUserInput{
			Email: "user@example.com",
		})

		if !errors.Is(err, ErrLoginPasswordRequired) {
			t.Fatalf("expected error %v, got %v", ErrLoginPasswordRequired, err)
		}
	})

	t.Run("returns invalid credentials when user does not exist", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{}
		tokenGenerator := fakeTokenGenerator{token: "access-token"}

		useCase := NewLoginUserUseCase(repository, hasher, tokenGenerator)

		_, err := useCase.Execute(context.Background(), LoginUserInput{
			Email:    "missing@example.com",
			Password: "secret-password",
		})

		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected error %v, got %v", ErrInvalidCredentials, err)
		}
	})

	t.Run("returns invalid credentials when password compare fails", func(t *testing.T) {
		repository := newFakeUserRepository()
		repository.usersByEmail["user@example.com"] = domain.User{
			ID:           "user-1",
			Email:        "user@example.com",
			PasswordHash: "hashed-password",
		}

		hasher := fakePasswordHasher{
			compareErr: errors.New("password mismatch"),
		}
		tokenGenerator := fakeTokenGenerator{token: "access-token"}

		useCase := NewLoginUserUseCase(repository, hasher, tokenGenerator)

		_, err := useCase.Execute(context.Background(), LoginUserInput{
			Email:    "user@example.com",
			Password: "wrong-password",
		})

		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected error %v, got %v", ErrInvalidCredentials, err)
		}
	})

	t.Run("returns error when token generation fails", func(t *testing.T) {
		repository := newFakeUserRepository()
		repository.usersByEmail["user@example.com"] = domain.User{
			ID:           "user-1",
			Email:        "user@example.com",
			PasswordHash: "hashed-password",
		}

		tokenErr := errors.New("token generation failed")

		hasher := fakePasswordHasher{}
		tokenGenerator := fakeTokenGenerator{
			generateErr: tokenErr,
		}

		useCase := NewLoginUserUseCase(repository, hasher, tokenGenerator)

		_, err := useCase.Execute(context.Background(), LoginUserInput{
			Email:    "user@example.com",
			Password: "secret-password",
		})

		if !errors.Is(err, tokenErr) {
			t.Fatalf("expected error %v, got %v", tokenErr, err)
		}
	})
}
