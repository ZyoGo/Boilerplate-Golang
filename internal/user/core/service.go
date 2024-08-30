package core

import "context"

type Service interface {
	CreateUser(ctx context.Context, user User) (User, error)
	FindUsers(ctx context.Context, filter FindUserFilter) ([]User, error)
	FindDetailUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, user User) error
}
