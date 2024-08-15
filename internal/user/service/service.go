package service

import (
	"github.com/ZyoGo/default-ddd-http/internal/user/modules/postgresql"
)

type service struct {
	userRepo postgresql.Repository
}

func New(userRepo postgresql.Repository) (*service, error) {
	return &service{
		userRepo,
	}, nil
}
