package service

import (
	"context"
	"fmt"

	user "github.com/ZyoGo/default-ddd-http/internal/user-v1/core"
)

type OptFunc func(u *UserService) error

func WithUserRepository(ur user.Repository) OptFunc {
	return func(u *UserService) (err error) {
		u.userRepo = ur
		return
	}
}

type UserService struct {
	userRepo user.Repository
}

func New(opts ...OptFunc) (user.Service, error) {
	us := &UserService{}

	for _, opt := range opts {
		if err := opt(us); err != nil {
			return nil, err
		}
	}

	if us.userRepo == nil {
		return nil, fmt.Errorf("user repository required")
	}

	return us, nil
}

func (svc *UserService) CreateUser(ctx context.Context, reqBody user.User) (user.User, error) {
	if err := reqBody.ValidatePassword(); err != nil {
		return user.User{}, err
	}

	return user.User{}, nil
}
