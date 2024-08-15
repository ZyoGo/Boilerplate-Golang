package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, body User) (User, error)
}

type User struct {
	ID       int64
	Email    string
	Password string
}
