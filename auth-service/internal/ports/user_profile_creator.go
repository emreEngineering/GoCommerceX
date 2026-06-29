package ports

import "context"

type UserProfile struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

type UserProfileCreator interface {
	Create(ctx context.Context, profile UserProfile) error
}
