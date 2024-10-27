package service

import (
	"context"

	"github.com/ZyoGo/default-ddd-http/internal/auth/core"
)

func (s *AuthService) SignIn(ctx context.Context, reqBody core.User) (core.Auth, error) {
	return core.Auth{}, nil
}
