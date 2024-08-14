package user

import (
	"fmt"

	"github.com/google/uuid"
)

type ErrInvalidUserId struct {
	UserId string
}

func (e ErrInvalidUserId) Error() string {
	return fmt.Sprintf("%s is invalid user id", e.UserId)
}

type User struct {
	Id         string
	TelegramUN string
}

func New(telegramUN string) User {
	return User{
		Id:         uuid.NewString(),
		TelegramUN: telegramUN,
	}
}

func ValidateUserId(userId string) error {
	if _, err := uuid.Parse(userId); err != nil {
		return ErrInvalidUserId{UserId: userId}
	}
	return nil
}
