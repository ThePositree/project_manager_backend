package user_managing_std

import (
	"context"
	"testing"

	"github.com/ThePositree/billing_manager/internal/usecase/mock"
	"github.com/stretchr/testify/require"
)

func TestUserManaging(t *testing.T) {
	mockRepo := &mock.MockUserRepository{}
	userManaging := New(mockRepo)

	ctx := context.Background()

	_, err := userManaging.Create(ctx, "test_telegram_un")
	require.NoError(t, err)

	_, err = userManaging.Create(ctx, "test_telegram_un")
	require.Error(t, err)

	users, err := userManaging.GetAll(ctx)
	require.NoError(t, err)

	require.Len(t, users, 1)

	_, err = userManaging.GetByTelegramUN(ctx, "test_telegram_un_2")
	require.Error(t, err)

	result, err := userManaging.GetByTelegramUN(ctx, "test_telegram_un")
	require.NoError(t, err)

	require.Equal(t, "test_telegram_un", result.TelegramUN)

	_, err = userManaging.GetById(ctx, "test_id_1")
	require.Error(t, err)

	result, err = userManaging.GetById(ctx, result.Id)
	require.NoError(t, err)

	require.Equal(t, "test_telegram_un", result.TelegramUN)

	_, err = userManaging.Delete(ctx, "test_id_1")
	require.Error(t, err)

	result, err = userManaging.Delete(ctx, result.Id)
	require.NoError(t, err)

	require.Equal(t, "test_telegram_un", result.TelegramUN)
}
