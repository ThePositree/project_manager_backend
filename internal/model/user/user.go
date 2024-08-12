package user

import "github.com/google/uuid"

type User struct {
	Id         string
	TelegramUN string
}

func New(telegramUN string) User {
	return User{
		Id:         uuid.New().String(),
		TelegramUN: telegramUN,
	}
}
