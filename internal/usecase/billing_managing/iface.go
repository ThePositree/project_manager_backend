package billing_managing

import (
	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/model/user"

	"github.com/ThePositree/billing_manager/internal/model/billing"
)

type BillingManaging interface {
	GetAllByUserId(userId string) ([]billing.Billing, error)
	GetById(id string) (billing.Billing, error)
	GetAll() ([]billing.Billing, error)
	Create(userId string) (billing.Billing, error)
	NextState(id string) (billing.Billing, error)
	PrevState(id string) (billing.Billing, error)
	SetBriefInfo(id string, username string) (billing.Billing, error)
}

type UserRepository interface {
	GetAll() ([]user.User, error)
	Get(id string) (user.User, error)
	Create(user user.User) (user.User, error)
	Delete(telegramUN string) (user.User, error)
}

type BillingRepository interface {
	GetAll() ([]model_billing.Billing, error)
	Get(id string) (model_billing.Billing, error)
	GetByUserId(userId string) ([]model_billing.Billing, error)
	Create(billing model_billing.Billing) (model_billing.Billing, error)
	Update(billing model_billing.Billing) (model_billing.Billing, error)
	Delete(id string) (model_billing.Billing, error)
}
