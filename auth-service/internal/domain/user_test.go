package domain

import (
	"errors"
	"testing"
	"time"
)

func TestUserValidate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name: "valid user",
			user: User{
				ID:           "user-1",
				Email:        "user@example.com",
				PasswordHash: "hashed-password",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "missing id",
			user: User{
				Email:        "user@example.com",
				PasswordHash: "hashed-password",
			},
			wantErr: ErrUserIDRequired,
		},
		{
			name: "missing email",
			user: User{
				ID:           "user-1",
				PasswordHash: "hashed-password",
			},
			wantErr: ErrUserEmailRequired,
		},
		{
			name: "missing password hash",
			user: User{
				ID:    "user-1",
				Email: "user@example.com",
			},
			wantErr: ErrPasswordHashRequired,
		},
		{
			name: "blank email",
			user: User{
				ID:           "user-1",
				Email:        "   ",
				PasswordHash: "hashed-password",
			},
			wantErr: ErrUserEmailRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
