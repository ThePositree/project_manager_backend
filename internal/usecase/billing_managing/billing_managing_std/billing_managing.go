package billing_managing_std

import (
	"fmt"

	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	billing_managing "github.com/ThePositree/billing_manager/internal/usecase/billing_managing"
)

var _ billing_managing.BillingManaging = billingManaging{}

type billingManaging struct {
	userRepo    billing_managing.UserRepository
	billingRepo billing_managing.BillingRepository
}

func (b billingManaging) Create(userId string) (model_billing.Billing, error) {
	user, err := b.userRepo.Get(userId)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting user by id from repository: %w", err)
	}

	billing, err := model_billing.New(user.Id)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("creating new billing from model: %w", err)
	}

	billing, err = b.billingRepo.Create(billing)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("creating new billing from repository: %w", err)
	}

	return billing, nil
}

func (b billingManaging) GetAll() ([]model_billing.Billing, error) {
	billings, err := b.billingRepo.GetAll()
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("getting all billings from repository: %w", err)
	}

	return billings, nil
}

func (b billingManaging) GetAllByUserId(userId string) ([]model_billing.Billing, error) {
	user, err := b.userRepo.Get(userId)
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("getting user by id from repository: %w", err)
	}

	billings, err := b.billingRepo.GetByUserId(user.Id)
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("getting billings by user id from repository: %w", err)
	}

	return billings, nil
}

func (b billingManaging) GetById(id string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(id)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}
	return billing, nil
}

func (b billingManaging) NextState(id string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(id)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}

	if err = billing.NextState(); err != nil {
		return model_billing.Billing{}, fmt.Errorf("billing next state: %w", err)
	}

	return billing, nil
}

func (b billingManaging) PrevState(id string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(id)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}

	if err = billing.PrevState(); err != nil {
		return model_billing.Billing{}, fmt.Errorf("billing next state: %w", err)
	}

	return billing, nil
}

func (b billingManaging) SetBriefInfo(id string, username string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(id)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}

	_, err = billing.SetBriefInfo(username)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("set brief info: %w", err)
	}

	return billing, nil
}

func New(userRepo billing_managing.UserRepository, billingRepo billing_managing.BillingRepository) billing_managing.BillingManaging {
	return billingManaging{
		userRepo:    userRepo,
		billingRepo: billingRepo,
	}
}
