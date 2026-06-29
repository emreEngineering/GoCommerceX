package grpc

import (
	"context"
	"errors"
	"log"

	"GoCommerceX/user-service/internal/application"
	"GoCommerceX/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	createUserUseCase       *application.CreateUserUseCase
	getUserUseCase          *application.GetUserUseCase
	getUserByEmailUseCase   *application.GetUserByEmailUseCase
	updateUserUseCase       *application.UpdateUserUseCase
	deleteUserUseCase       *application.DeleteUserUseCase
}

func NewUserHandler(
	createUserUseCase *application.CreateUserUseCase,
	getUserUseCase *application.GetUserUseCase,
	getUserByEmailUseCase *application.GetUserByEmailUseCase,
	updateUserUseCase *application.UpdateUserUseCase,
	deleteUserUseCase *application.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:       createUserUseCase,
		getUserUseCase:          getUserUseCase,
		getUserByEmailUseCase:   getUserByEmailUseCase,
		updateUserUseCase:       updateUserUseCase,
		deleteUserUseCase:       deleteUserUseCase,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	input := application.CreateUserInput{
		Email:     req.GetEmail(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
	}

	output, err := h.createUserUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCreateUserEmailRequired),
			errors.Is(err, application.ErrCreateUserFirstNameRequired),
			errors.Is(err, application.ErrCreateUserLastNameRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			log.Printf("CreateUser error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userv1.CreateUserResponse{
		User: &userv1.User{
			Id:        output.User.ID,
			Email:     output.User.Email,
			FirstName: output.User.FirstName,
			LastName:  output.User.LastName,
			Phone:     output.User.Phone,
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	input := application.GetUserInput{ID: req.GetId()}
	output, err := h.getUserUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetUser error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:        output.User.ID,
			Email:     output.User.Email,
			FirstName: output.User.FirstName,
			LastName:  output.User.LastName,
			Phone:     output.User.Phone,
		},
	}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *userv1.GetUserByEmailRequest) (*userv1.GetUserByEmailResponse, error) {
	input := application.GetUserByEmailInput{Email: req.GetEmail()}
	output, err := h.getUserByEmailUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetUserEmailRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetUserByEmail error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userv1.GetUserByEmailResponse{
		User: &userv1.User{
			Id:        output.User.ID,
			Email:     output.User.Email,
			FirstName: output.User.FirstName,
			LastName:  output.User.LastName,
			Phone:     output.User.Phone,
		},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	input := application.UpdateUserInput{
		ID:        req.GetId(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
	}
	output, err := h.updateUserUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrUpdateUserIDRequired),
			errors.Is(err, application.ErrUpdateUserFirstNameRequired),
			errors.Is(err, application.ErrUpdateUserLastNameRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("UpdateUser error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:        output.User.ID,
			Email:     output.User.Email,
			FirstName: output.User.FirstName,
			LastName:  output.User.LastName,
			Phone:     output.User.Phone,
		},
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	input := application.DeleteUserInput{ID: req.GetId()}
	output, err := h.deleteUserUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrDeleteUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("DeleteUser error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userv1.DeleteUserResponse{Success: output.Success}, nil
}
