package billing

import (
	"fmt"

	"github.com/ThePositree/billing_manager/internal/model/user"
	"github.com/google/uuid"
)

type ErrNextCompletedState struct{}

func (e ErrNextCompletedState) Error() string {
	return "impossible to next the state from completed state"
}

type ErrPrevPendingState struct{}

func (e ErrPrevPendingState) Error() string {
	return "impossible to prev the state from pending state"
}

type BriefInfo struct {
	Username string
}

type Billing struct {
	Id        string
	UserId    string
	_state    State
	_username string
}

func New(userId string) (Billing, error) {
	err := user.ValidateUserId(userId)
	if err != nil {
		return Billing{}, err
	}

	return Billing{
		Id:     uuid.New().String(),
		_state: StatePending,
		UserId: userId,
	}, nil
}

func (b *Billing) NextState() error {
	switch b._state {
	case StatePending:
		b._state = StateDesign
		return nil
	case StateDesign:
		b._state = StateLayout
		return nil
	case StateLayout:
		b._state = StateCompleted
		return nil
	case StateCompleted:
		return ErrNextCompletedState{}
	}
	return fmt.Errorf("%s is %w", b._state, ErrInvalidState)
}

func (b *Billing) PrevState() error {
	switch b._state {
	case StatePending:
		return ErrPrevPendingState{}
	case StateDesign:
		b._state = StatePending
		return nil
	case StateLayout:
		b._state = StateDesign
		return nil
	case StateCompleted:
		b._state = StateLayout
		return nil
	}
	return fmt.Errorf("%s is %w", b._state, ErrInvalidState)
}

func (b *Billing) GetState() State {
	return b._state
}

func (b *Billing) GetBriefInfo() BriefInfo {
	return BriefInfo{
		Username: b._username,
	}
}

func (b *Billing) SetBriefInfo(username string) (BriefInfo, error) {
	b._username = username
	return BriefInfo{Username: b._username}, nil
}

// ENUM(
// pending
// design
// layout
// completed
// )
type State string

type ErrInvalidBillingId struct {
	BillingId string
}

func (e ErrInvalidBillingId) Error() string {
	return fmt.Sprintf("%s is invalid user id", e.BillingId)
}

func ValidateBillingId(billingId string) error {
	if _, err := uuid.Parse(billingId); err != nil {
		return ErrInvalidBillingId{BillingId: billingId}
	}
	return nil
}

type DTO interface {
	GetId() string
	GetUserId() string
	GetState() string
	GetUsername() string
}

func ToModelFromDTO(dto DTO) (Billing, error) {
	dtoState := dto.GetState()
	state, err := ParseState(dtoState)
	if err != nil {
		return Billing{}, fmt.Errorf("%s is %w", dtoState, ErrInvalidState)
	}
	userId := dto.GetUserId()
	err = user.ValidateUserId(userId)
	if err != nil {
		return Billing{}, err
	}
	id := dto.GetId()
	err = ValidateBillingId(id)
	if err != nil {
		return Billing{}, err
	}
	return Billing{
		Id:        id,
		UserId:    userId,
		_state:    state,
		_username: dto.GetUsername(),
	}, nil
}
