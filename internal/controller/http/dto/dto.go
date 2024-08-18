package dto

import (
	"github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/model/user"
)

type User struct {
	Id         string `json:"id"`
	TelegramUN string `json:"telegram_username"`
}

type CreateUserInfo struct {
	TelegramUN string `json:"telegram_username"`
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

type Billing struct {
	Id       string `json:"id"`
	UserId   string `json:"user_id"`
	State    string `json:"state"`
	Username string `json:"username"`
}

type CreateBillingInfo struct {
	UserId string `json:"user_id"`
}

func (u Billing) GetUsername() string {
	return u.Username
}

func (u Billing) ToModel() (billing.Billing, error) {
	return billing.ToModelFromDTO(u)
}

func (u Billing) GetId() string {
	return u.Id
}

func (u Billing) GetUserId() string {
	return u.UserId
}

func (u Billing) GetState() string {
	return u.State
}

func NewBillingDTOFromModel(billing billing.Billing) Billing {
	return Billing{
		Id:       billing.Id,
		UserId:   billing.UserId,
		State:    billing.GetState().String(),
		Username: billing.GetBriefInfo().Username,
	}
}

type BriefInfo struct {
	Username string `json:"username"`
}
