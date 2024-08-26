package mock

import (
	"context"
	"fmt"
	"slices"

	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/usecase"
)

var _ usecase.BillingRepository = &MockBillingRepository{}

type MockBillingRepository struct {
	billings []model_billing.Billing
}

func (m *MockBillingRepository) GetByUserId(ctx context.Context, userId string) ([]model_billing.Billing, error) {
	var result []model_billing.Billing
	for _, billing := range m.billings {
		if billing.UserId == userId {
			result = append(result, billing)
		}
	}
	return result, nil
}

func (m *MockBillingRepository) Update(ctx context.Context, billingForUpdate model_billing.Billing) (model_billing.Billing, error) {
	index := -1
	result := model_billing.Billing{}
	for i, billing := range m.billings {
		if billingForUpdate.Id == billing.Id {
			index = i
			billing = billingForUpdate
			result = billing
		}
	}
	if index >= 0 {
		return result, nil
	}
	return result, ErrNoData
}

func (m *MockBillingRepository) Get(ctx context.Context, id string) (model_billing.Billing, error) {
	for _, billing := range m.billings {
		if billing.Id == id {
			return billing, nil
		}
	}
	return model_billing.Billing{}, ErrNoData
}

func (m *MockBillingRepository) GetNoDataError() error {
	return ErrNoData
}

func (m *MockBillingRepository) Create(ctx context.Context, billing model_billing.Billing) (model_billing.Billing, error) {
	m.billings = append(m.billings, billing)
	fmt.Println(m.billings)
	return billing, nil
}

func (m *MockBillingRepository) GetAll(ctx context.Context) ([]model_billing.Billing, error) {
	return m.billings, nil
}

func (m *MockBillingRepository) Delete(ctx context.Context, id string) (model_billing.Billing, error) {
	result := model_billing.Billing{}
	var index int = -1
	for i, billing := range m.billings {
		if billing.Id == id {
			result = billing
			index = i
		}
	}

	if index >= 0 {
		m.billings = slices.Delete(m.billings, index, index+1)
		return result, nil
	}

	return result, ErrNoData
}
