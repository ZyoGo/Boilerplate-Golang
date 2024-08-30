package http

import (
	"github.com/ZyoGo/default-ddd-http/internal/user-v1/core"
	"github.com/ZyoGo/default-ddd-http/internal/user-v1/infrastructure/http/v1/request"
)

func CreateUserDTO(req *request.CreateUser) (res core.User) {
	res = core.User{
		Email:    req.Email,
		Password: req.Password,
	}

	return
}
