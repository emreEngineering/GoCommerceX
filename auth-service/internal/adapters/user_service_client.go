package adapters

import (
	"context"
	"errors"

	"GoCommerceX/auth-service/internal/application"
	"GoCommerceX/auth-service/internal/ports"
	userv1 "GoCommerceX/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceClient struct {
	client userv1.UserServiceClient
}

func NewUserServiceClient(client userv1.UserServiceClient) *UserServiceClient {
	return &UserServiceClient{client: client}
}

func (c *UserServiceClient) Create(ctx context.Context, profile ports.UserProfile) error {
	_, err := c.client.CreateUser(ctx, &userv1.CreateUserRequest{
		Id:        profile.ID,
		Email:     profile.Email,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Phone:     profile.Phone,
	})
	if err == nil {
		return nil
	}

	var appErr error
	switch status.Code(err) {
	case codes.AlreadyExists:
		appErr = application.ErrUserProfileAlreadyExists
	case codes.InvalidArgument, codes.FailedPrecondition:
		appErr = application.ErrUserProfileCreationFailed
	default:
		appErr = application.ErrUserProfileCreationFailed
	}

	return errors.Join(appErr, err)
}
