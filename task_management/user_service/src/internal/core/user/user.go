package user

import (
	"context"
	"task-management/user-service/src/internal/core/session"
)

type Repository interface {
	RegisterUser(ctx context.Context, user *session.RegisterResponse) error

	Login(ctx context.Context, user *session.RegisterResponse) (*session.RegisterResponse, error)

	IsEmailOsUserNameTaken(ctx context.Context, email, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*session.RegisterResponse, error)
	GetUserByID(ctx context.Context, uid int) (*session.RegisterResponse, error)
	UpdateUser(ctx context.Context, user *session.RegisterResponse) error
}

