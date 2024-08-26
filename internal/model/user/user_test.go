package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockDTO struct {
	Id         string
	TelegramUN string
}

func (m mockDTO) GetId() string {
	return m.Id
}

func (m mockDTO) GetTelegramUN() string {
	return m.TelegramUN
}

func TestToModelFromDTO(t *testing.T) {
	userId := uuid.NewString()

	test := mockDTO{
		Id:         userId,
		TelegramUN: "test",
	}

	user, err := ToModelFromDTO(test)
	require.NoError(t, err)

	require.Equal(t, User{
		Id:         userId,
		TelegramUN: "test",
	}, user)
}

func TestValidateUserId(t *testing.T) {
	err := ValidateUserId("test")
	require.Error(t, err)

	userId := uuid.NewString()
	err = ValidateUserId(userId)
	require.NoError(t, err)
}
