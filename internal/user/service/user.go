package service

import (
	"context"
	"time"

	user "github.com/ZyoGo/default-ddd-http/internal/user/core"
)

func (svc *UserService) CreateUser(ctx context.Context, reqBody user.User) (user.User, error) {
	if err := reqBody.ValidatePassword(); err != nil {
		return user.User{}, err
	}

	userData, err := svc.userRepo.FindUserByEmail(ctx, reqBody.Email)
	if err != nil {
		return user.User{}, err
	}

	// check if user already exist
	if userData.Email != "" {
		return user.User{}, user.ErrUserAlreadyExist
	}

	// generate id using contract (ulid)
	reqBody.GenerateID(svc.ID.Generate())
	if err := svc.userRepo.InsertUser(ctx, reqBody); err != nil {
		return user.User{}, err
	}

	return reqBody, nil
}

func (svc *UserService) FindUsers(ctx context.Context, filter user.FindUserFilter) ([]user.User, error) {
	users, err := svc.userRepo.FindUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (svc *UserService) FindDetailUser(ctx context.Context, id string) (user.User, error) {
	result, err := svc.userRepo.FindUserByID(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	return result, nil
}

func (svc *UserService) UpdateUser(ctx context.Context, reqBody user.User) error {
	if err := reqBody.ValidatePassword(); err != nil {
		return err
	}

	userData, err := svc.userRepo.FindUserByID(ctx, reqBody.ID)
	if err != nil {
		return err
	}

	userUpdate := user.User{
		ID:        userData.ID,
		Email:     userData.Email,
		Password:  reqBody.Password,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}

	if err := svc.userRepo.UpdateUser(ctx, userUpdate); err != nil {
		return err
	}

	return nil
}
