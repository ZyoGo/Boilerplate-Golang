package service

import (
	"fmt"

	"github.com/ZyoGo/default-ddd-http/internal/user/core"
	user "github.com/ZyoGo/default-ddd-http/internal/user/core"
)

type OptFunc func(u *UserService) error

func WithUserRepository(ur user.Repository) OptFunc {
	return func(u *UserService) (err error) {
		u.userRepo = ur
		return
	}
}

func WithIDGenerator(id user.ID) OptFunc {
	return func(u *UserService) (err error) {
		u.ID = id
		return
	}
}

func WithHash(hash core.Hash) OptFunc {
	return func(u *UserService) (err error) {
		u.Hash = hash
		return
	}
}

type UserService struct {
	userRepo user.Repository
	ID       user.ID
	Hash     user.Hash
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

	if us.ID == nil {
		return nil, fmt.Errorf("id generator required")
	}

	if us.Hash == nil {
		return nil, fmt.Errorf("hash package required")
	}

	return us, nil
}
