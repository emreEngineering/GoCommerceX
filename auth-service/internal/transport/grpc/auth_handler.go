package grpc

import (
	"context"
	"errors"

	"GoCommerceX/auth-service/internal/application"
	"GoCommerceX/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	registerUseCase *application.RegisterUserUseCase
	loginUseCase    *application.LoginUserUseCase
}

func NewAuthHandler(registerUseCase *application.RegisterUserUseCase, loginUseCase *application.LoginUserUseCase) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	input := application.RegisterUserInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	output, err := h.registerUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrRegisterEmailRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrRegisterPasswordRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authv1.RegisterResponse{
		UserId: output.UserID,
		Email:  output.Email,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	input := application.LoginUserInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	output, err := h.loginUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrLoginEmailRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrLoginPasswordRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authv1.LoginResponse{
		AccessToken: output.AccessToken,
	}, nil
}
