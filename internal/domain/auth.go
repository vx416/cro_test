package domain

import (
	"context"
	"cro_test/internal/model"
)

type AuthRepositorier interface {
	GetUser(ctx context.Context, opts model.GetUserOpts) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type AuthServicer interface {
	SignUp(ctx context.Context, email, password string) (model.User, error)
	Login(ctx context.Context, email, password string) (model.User, error)
}
