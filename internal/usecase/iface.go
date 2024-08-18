package usecase

import (
	"context"

	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/model/user"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]user.User, error)
	GetByTelegramUN(ctx context.Context, telegramUN string) (user.User, error)
	Get(ctx context.Context, id string) (user.User, error)
	Create(ctx context.Context, user user.User) (user.User, error)
	Delete(ctx context.Context, id string) (user.User, error)
	GetNoDataError() error
}

type BillingRepository interface {
	GetAll(ctx context.Context) ([]model_billing.Billing, error)
	Get(ctx context.Context, id string) (model_billing.Billing, error)
	GetByUserId(ctx context.Context, userId string) ([]model_billing.Billing, error)
	Create(ctx context.Context, billing model_billing.Billing) (model_billing.Billing, error)
	Update(ctx context.Context, billing model_billing.Billing) (model_billing.Billing, error)
	Delete(ctx context.Context, id string) (model_billing.Billing, error)
	GetNoDataError() error
}
