package billing_managing

import (
	"context"
	"errors"

	"github.com/ThePositree/billing_manager/internal/model/billing"
)

var (
	ErrNoUser    = errors.New("no data")
	ErrNoBilling = errors.New("no data")
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
