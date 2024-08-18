package dto

import "github.com/ThePositree/billing_manager/internal/model/user"

type User struct {
	Id         string `bson:"_id"`
	TelegramUN string `bson:"telegram_username"`
}

func (u User) GetId() string {
	return u.Id
}

func (u User) GetTelegramUN() string {
	return u.TelegramUN
}

func (u User) ToModel() (user.User, error) {
	return user.ToModelFromDTO(u)
}

func NewUserDTOFromModel(user user.User) User {
	return User{
		Id:         user.Id,
		TelegramUN: user.TelegramUN,
	}
}
