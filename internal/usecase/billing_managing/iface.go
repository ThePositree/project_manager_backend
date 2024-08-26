package billing_managing

import (
	"context"
	"errors"

	"github.com/ThePositree/billing_manager/internal/model/billing"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrBillingNotFound = errors.New("billing not found")
)

type BillingManaging interface {
	GetAllByUserId(ctx context.Context, userId string) ([]billing.Billing, error)
	GetById(ctx context.Context, id string) (billing.Billing, error)
	GetAll(ctx context.Context) ([]billing.Billing, error)
	Create(ctx context.Context, userId string) (billing.Billing, error)
	NextState(ctx context.Context, id string) (billing.Billing, error)
	PrevState(ctx context.Context, id string) (billing.Billing, error)
	SetBriefInfo(ctx context.Context, id string, username string) (billing.Billing, error)
}
