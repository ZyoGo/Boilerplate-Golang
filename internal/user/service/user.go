package service

import (
	"context"
	"errors"

	"github.com/ZyoGo/default-ddd-http/internal/user"
)

func (s *service) CreateUser(ctx context.Context, reqBody user.User) (user.User, error) {
	ctx, err := s.userRepo.TransactionContext(ctx)
	if err != nil {
		return user.User{}, err
	}
	defer s.userRepo.Rollback(ctx)

	if _, err := s.userRepo.FindUserByEmail(ctx, reqBody.Email); err != nil {
		return user.User{}, errors.New("user already exists")
	}

	if err := s.userRepo.InsertUser(ctx, reqBody); err != nil {
		return user.User{}, err
	}

	if err := s.userRepo.UpdateUser(ctx, reqBody); err != nil {
		return user.User{}, err
	}

	if err := s.userRepo.Commit(ctx); err != nil {
		return user.User{}, err
	}

	return reqBody, nil
}
