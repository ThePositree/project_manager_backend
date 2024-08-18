package user_managing

import (
	"context"
	"errors"

	"github.com/ThePositree/billing_manager/internal/model/user"
)

var ErrExistingUser = errors.New("user is existing")

type UserManaging interface {
	GetByTelegramUN(ctx context.Context, telegramUN string) (user.User, error)
	GetById(ctx context.Context, id string) (user.User, error)
	GetAll(ctx context.Context) ([]user.User, error)
	Create(ctx context.Context, telegramUN string) (user.User, error)
	Delete(ctx context.Context, id string) (user.User, error)
}
