package service

import (
	"context"
	"reset/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) dto.UserResponse
	LoginUser(ctx context.Context, req dto.LoginRequest) (string, error)
	FindByNRA(ctx context.Context, nra string) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, request dto.ChangePasswordRequest) error
}