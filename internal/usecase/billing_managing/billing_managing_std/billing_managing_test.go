package billing_managing_std

import (
	"context"
	"testing"

	"github.com/ThePositree/billing_manager/internal/model/user"
	"github.com/ThePositree/billing_manager/internal/usecase/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBillingManaging(t *testing.T) {
	mockUserRepo := &mock.MockUserRepository{}
	mockBillingRepo := &mock.MockBillingRepository{}
	billingManaging := New(mockUserRepo, mockBillingRepo)

	ctx := context.Background()

	userId := uuid.NewString()

	_, err := billingManaging.Create(ctx, userId)
	require.Error(t, err)

	_, err = mockUserRepo.Create(ctx, user.User{Id: userId, TelegramUN: "test"})
	require.NoError(t, err)

	_, err = billingManaging.Create(ctx, userId)
	require.NoError(t, err)

	billings, err := billingManaging.GetAll(ctx)
	require.NoError(t, err)

	require.Len(t, billings, 1)
}
