package core

import "context"

type Service interface {
	SignIn(ctx context.Context, user User) (Auth, error)
}
