package user_managing_std

import (
	"context"
	"errors"
	"fmt"

	model_user "github.com/ThePositree/billing_manager/internal/model/user"
	"github.com/ThePositree/billing_manager/internal/usecase"
	"github.com/ThePositree/billing_manager/internal/usecase/user_managing"
)

var _ user_managing.UserManaging = userManaging{}

type userManaging struct {
	userRepo usecase.UserRepository
}

func (u userManaging) Create(ctx context.Context, telegramUN string) (model_user.User, error) {
	var user model_user.User
	_, err := u.userRepo.GetByTelegramUN(ctx, telegramUN)
	if err == nil {
		return model_user.User{}, user_managing.ErrExistingUser
	}
	if errors.Is(u.userRepo.GetNoDataError(), err) {
		user = model_user.New(telegramUN)
	} else {
		return model_user.User{}, fmt.Errorf("getting user by id from repository: %w", err)
	}

	user, err = u.userRepo.Create(ctx, user)
	if err != nil {
		return model_user.User{}, fmt.Errorf("creating new user from repository: %w", err)
	}

	return user, nil
}

func (u userManaging) GetAll(ctx context.Context) ([]model_user.User, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return []model_user.User{}, fmt.Errorf("getting all users from repository: %w", err)
	}

	return users, nil
}

func (u userManaging) GetByTelegramUN(ctx context.Context, telegramUN string) (model_user.User, error) {
	user, err := u.userRepo.GetByTelegramUN(ctx, telegramUN)
	if err != nil {
		return model_user.User{}, fmt.Errorf("getting user by id from repository: %w", err)
	}
	return user, nil
}

func (u userManaging) GetById(ctx context.Context, id string) (model_user.User, error) {
	user, err := u.userRepo.Get(ctx, id)
	if err != nil {
		return model_user.User{}, fmt.Errorf("getting user by id from repository: %w", err)
	}
	return user, nil
}

func (u userManaging) Delete(ctx context.Context, id string) (model_user.User, error) {
	user, err := u.userRepo.Delete(ctx, id)
	if err != nil {
		return model_user.User{}, fmt.Errorf("deleting user from repository: %w", err)
	}
	return user, nil
}

func New(userRepo usecase.UserRepository) userManaging {
	return userManaging{
		userRepo: userRepo,
	}
}
