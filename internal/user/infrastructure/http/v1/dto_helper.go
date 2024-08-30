package http

import (
	"github.com/ZyoGo/default-ddd-http/internal/user/core"
	"github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http/v1/request"
)

func CreateUserDTO(req *request.CreateUser) (res core.User) {
	res = core.User{
		Email:    req.Email,
		Password: req.Password,
	}

	return
}
