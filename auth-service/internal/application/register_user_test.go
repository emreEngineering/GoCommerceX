package application

import (
	"context"
	"errors"
	"testing"

	"GoCommerceX/auth-service/internal/domain"
	"GoCommerceX/auth-service/internal/ports"
)

type fakeUserRepository struct {
	usersByEmail  map[string]domain.User
	savedUser     domain.User
	deletedUserID string
	saveErr       error
	findErr       error
	deleteErr     error
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

func (r *fakeUserRepository) Delete(ctx context.Context, id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}

	r.deletedUserID = id
	delete(r.usersByEmail, id)
	return nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if r.findErr != nil {
		return domain.User{}, r.findErr
	}

	user, ok := r.usersByEmail[email]
	if !ok {
		return domain.User{}, ports.ErrUserNotFound
	}

	return user, nil
}

type fakePasswordHasher struct {
	hashValue  string
	hashErr    error
	compareErr error
}

func (h fakePasswordHasher) Hash(ctx context.Context, plainPassword string) (string, error) {
	if h.hashErr != nil {
		return "", h.hashErr
	}

	return h.hashValue, nil
}

func (h fakePasswordHasher) Compare(ctx context.Context, plainPassword string, passwordHash string) error {
	return h.compareErr
}

type fakeUserProfileCreator struct {
	createErr error
	created   ports.UserProfile
	called    bool
}

func (c *fakeUserProfileCreator) Create(ctx context.Context, profile ports.UserProfile) error {
	c.called = true
	c.created = profile
	return c.createErr
}

func TestRegisterUserUseCaseExecute(t *testing.T) {
	t.Run("registers user successfully", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{
			hashValue: "hashed-password",
		}
		profileCreator := &fakeUserProfileCreator{}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		output, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Email:     " user@example.com ",
			Password:  " secret-password ",
			FirstName: " Emre ",
			LastName:  " Developer ",
			Phone:     " 555-0100 ",
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

		if !profileCreator.called {
			t.Fatalf("expected user profile creator to be called")
		}

		if profileCreator.created.ID != "user-1" {
			t.Fatalf("expected profile id user-1, got %s", profileCreator.created.ID)
		}
	})

	t.Run("returns error when email is missing", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{hashValue: "hashed-password"}
		profileCreator := &fakeUserProfileCreator{}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Password:  "secret-password",
			FirstName: "Emre",
			LastName:  "Developer",
		})

		if !errors.Is(err, ErrRegisterEmailRequired) {
			t.Fatalf("expected error %v, got %v", ErrRegisterEmailRequired, err)
		}
	})

	t.Run("returns error when password is missing", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{hashValue: "hashed-password"}
		profileCreator := &fakeUserProfileCreator{}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Email:     "user@example.com",
			FirstName: "Emre",
			LastName:  "Developer",
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
		profileCreator := &fakeUserProfileCreator{}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Email:     "user@example.com",
			Password:  "secret-password",
			FirstName: "Emre",
			LastName:  "Developer",
		})

		if !errors.Is(err, ErrUserAlreadyExists) {
			t.Fatalf("expected error %v, got %v", ErrUserAlreadyExists, err)
		}
	})

	t.Run("returns error when password hashing fails", func(t *testing.T) {
		repository := newFakeUserRepository()
		hashErr := errors.New("hash failed")
		profileCreator := &fakeUserProfileCreator{}

		hasher := fakePasswordHasher{
			hashErr: hashErr,
		}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Email:     "user@example.com",
			Password:  "secret-password",
			FirstName: "Emre",
			LastName:  "Developer",
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
		profileCreator := &fakeUserProfileCreator{}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Email:     "user@example.com",
			Password:  "secret-password",
			FirstName: "Emre",
			LastName:  "Developer",
		})

		if !errors.Is(err, saveErr) {
			t.Fatalf("expected error %v, got %v", saveErr, err)
		}
	})

	t.Run("returns error when profile creation fails", func(t *testing.T) {
		repository := newFakeUserRepository()
		hasher := fakePasswordHasher{hashValue: "hashed-password"}
		profileCreator := &fakeUserProfileCreator{createErr: errors.New("profile create failed")}

		useCase := NewRegisterUserUseCase(repository, hasher, profileCreator)

		_, err := useCase.Execute(context.Background(), RegisterUserInput{
			ID:        "user-1",
			Email:     "user@example.com",
			Password:  "secret-password",
			FirstName: "Emre",
			LastName:  "Developer",
		})

		if !errors.Is(err, ErrUserProfileCreationFailed) {
			t.Fatalf("expected error %v, got %v", ErrUserProfileCreationFailed, err)
		}
	})
}
