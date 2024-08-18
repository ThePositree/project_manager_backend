package dto

import (
	"github.com/ThePositree/billing_manager/internal/model/billing"
)

type Billing struct {
	Id       string `bson:"_id"`
	UserId   string `bson:"user_id"`
	State    string `bson:"state"`
	Username string `bson:"username"`
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
