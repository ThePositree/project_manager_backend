package mock

import (
	"context"
	"errors"
	"fmt"
	"slices"

	model_user "github.com/ThePositree/billing_manager/internal/model/user"
	"github.com/ThePositree/billing_manager/internal/usecase"
)

var _ usecase.UserRepository = &MockUserRepository{}

type MockUserRepository struct {
	users []model_user.User
}

var ErrNoData error = errors.New("no data")

func (m *MockUserRepository) Get(ctx context.Context, id string) (model_user.User, error) {
	for _, user := range m.users {
		if user.Id == id {
			return user, nil
		}
	}
	return model_user.User{}, ErrNoData
}

func (m *MockUserRepository) GetNoDataError() error {
	return ErrNoData
}

func (m *MockUserRepository) GetByTelegramUN(ctx context.Context, telegramUN string) (model_user.User, error) {
	for _, user := range m.users {
		if user.TelegramUN == telegramUN {
			return user, nil
		}
	}
	return model_user.User{}, ErrNoData
}

func (m *MockUserRepository) Create(ctx context.Context, user model_user.User) (model_user.User, error) {
	m.users = append(m.users, user)
	fmt.Println(m.users)
	return user, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]model_user.User, error) {
	return m.users, nil
}

func (m *MockUserRepository) GetById(ctx context.Context, id string) (model_user.User, error) {
	for _, user := range m.users {
		if user.Id == id {
			return user, nil
		}
	}
	return model_user.User{}, ErrNoData
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) (model_user.User, error) {
	result := model_user.User{}
	var index int = -1
	for i, user := range m.users {
		if user.Id == id {
			result = user
			index = i
		}
	}

	if index >= 0 {
		m.users = slices.Delete(m.users, index, index+1)
		return result, nil
	}

	return result, ErrNoData
}
