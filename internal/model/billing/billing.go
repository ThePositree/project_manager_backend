package billing

import (
	"fmt"

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

type Billing struct {
	Id     string
	UserId string
	state  State
}

func New(userId string) (Billing, error) {
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return Billing{}, fmt.Errorf("parsing user id: %w", err)
	}

	return Billing{
		Id:     uuid.New().String(),
		state:  StatePending,
		UserId: userUuid.String(),
	}, nil
}

func (b Billing) NextState() error {
	switch b.state {
	case StatePending:
		b.state = StateDesign
		return nil
	case StateDesign:
		b.state = StateLayout
		return nil
	case StateLayout:
		b.state = StateCompleted
		return nil
	case StateCompleted:
		return ErrNextCompletedState{}
	}
	return fmt.Errorf("%s is %w", b.state, ErrInvalidState)
}

func (b Billing) PrevState() error {
	switch b.state {
	case StatePending:
		return ErrPrevPendingState{}
	case StateDesign:
		b.state = StatePending
		return nil
	case StateLayout:
		b.state = StateDesign
		return nil
	case StateCompleted:
		b.state = StateLayout
	}
	return fmt.Errorf("%s is %w", b.state, ErrInvalidState)
}

func (b Billing) GetState() State {
	return b.state
}

// ENUM(
// pending
// design
// layout
// completed
// )
type State string
