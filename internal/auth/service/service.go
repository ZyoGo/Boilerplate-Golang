package service

import (
	"fmt"

	auth "github.com/ZyoGo/default-ddd-http/internal/auth/core"
)

type OptFunc func(a *AuthService) error

func WithAuthRepository(ar auth.Repository) OptFunc {
	return func(a *AuthService) (err error) {
		a.authRepo = ar
		return
	}
}

func WithIDGenerator(id auth.ID) OptFunc {
	return func(a *AuthService) (err error) {
		a.ID = id
		return
	}
}

type AuthService struct {
	authRepo auth.Repository
	ID       auth.ID
}

func New(opts ...OptFunc) (auth.Service, error) {
	as := &AuthService{}

	for _, opt := range opts {
		if err := opt(as); err != nil {
			return nil, err
		}
	}

	if as.authRepo == nil {
		return nil, fmt.Errorf("auth repository required")
	}

	if as.ID == nil {
		return nil, fmt.Errorf("id generator required")
	}

	return as, nil
}
