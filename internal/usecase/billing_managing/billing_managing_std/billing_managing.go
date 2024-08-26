package billing_managing_std

import (
	"context"
	"errors"
	"fmt"

	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/usecase"
	"github.com/ThePositree/billing_manager/internal/usecase/billing_managing"
)

var _ billing_managing.BillingManaging = billingManaging{}

type billingManaging struct {
	userRepo    usecase.UserRepository
	billingRepo usecase.BillingRepository
}

func (b billingManaging) Create(ctx context.Context, userId string) (model_billing.Billing, error) {
	user, err := b.userRepo.Get(ctx, userId)
	if errors.Is(b.userRepo.GetNoDataError(), err) {
		return model_billing.Billing{}, billing_managing.ErrUserNotFound
	}
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting user by id from repository: %w", err)
	}

	billing, err := model_billing.New(user.Id)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("creating new billing from model: %w", err)
	}

	billing, err = b.billingRepo.Create(ctx, billing)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("creating new billing from repository: %w", err)
	}

	return billing, nil
}

func (b billingManaging) GetAll(ctx context.Context) ([]model_billing.Billing, error) {
	billings, err := b.billingRepo.GetAll(ctx)
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("getting all billings from repository: %w", err)
	}

	return billings, nil
}

func (b billingManaging) GetAllByUserId(ctx context.Context, userId string) ([]model_billing.Billing, error) {
	user, err := b.userRepo.Get(ctx, userId)
	if errors.Is(b.userRepo.GetNoDataError(), err) {
		return []model_billing.Billing{}, billing_managing.ErrUserNotFound
	}
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("getting user by id from repository: %w", err)
	}

	billings, err := b.billingRepo.GetByUserId(ctx, user.Id)
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("getting billings by user id from repository: %w", err)
	}

	return billings, nil
}

func (b billingManaging) GetById(ctx context.Context, id string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(ctx, id)
	if errors.Is(b.billingRepo.GetNoDataError(), err) {
		return model_billing.Billing{}, billing_managing.ErrBillingNotFound
	}
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}
	return billing, nil
}

func (b billingManaging) NextState(ctx context.Context, id string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(ctx, id)
	if errors.Is(b.billingRepo.GetNoDataError(), err) {
		return model_billing.Billing{}, billing_managing.ErrBillingNotFound
	}
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}

	if err = billing.NextState(); err != nil {
		return model_billing.Billing{}, err
	}

	billing, err = b.billingRepo.Update(ctx, billing)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("updating billing in repository: %w", err)
	}

	return billing, nil
}

func (b billingManaging) PrevState(ctx context.Context, id string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(ctx, id)
	if errors.Is(b.billingRepo.GetNoDataError(), err) {
		return model_billing.Billing{}, billing_managing.ErrBillingNotFound
	}
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}

	if err = billing.PrevState(); err != nil {
		return model_billing.Billing{}, err
	}

	billing, err = b.billingRepo.Update(ctx, billing)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("updating billing in repository: %w", err)
	}

	return billing, nil
}

func (b billingManaging) SetBriefInfo(ctx context.Context, id string, username string) (model_billing.Billing, error) {
	billing, err := b.billingRepo.Get(ctx, id)
	if errors.Is(b.billingRepo.GetNoDataError(), err) {
		return model_billing.Billing{}, billing_managing.ErrBillingNotFound
	}
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("getting billing by id from repository: %w", err)
	}

	_, err = billing.SetBriefInfo(username)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("set brief info: %w", err)
	}

	billing, err = b.billingRepo.Update(ctx, billing)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("updating billing in repository: %w", err)
	}

	return billing, nil
}

func New(userRepo usecase.UserRepository, billingRepo usecase.BillingRepository) billingManaging {
	return billingManaging{
		userRepo:    userRepo,
		billingRepo: billingRepo,
	}
}
