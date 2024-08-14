package billing

import (
	"testing"

	"github.com/ThePositree/billing_manager/internal/model/user"
	"github.com/stretchr/testify/assert"
)

func TestBilling(t *testing.T) {
	_, err := New("blablabla")
	assert.EqualError(t, err, (user.ErrInvalidUserId{UserId: "blablabla"}).Error())

	billing, err := New("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)

	err = billing.PrevState()
	assert.EqualError(t, err, (ErrPrevPendingState{}).Error())

	err = billing.NextState()
	assert.NoError(t, err)

	state := billing.GetState()
	assert.Equal(t, StateDesign, state)

	err = billing.PrevState()
	assert.NoError(t, err)

	err = billing.NextState()
	assert.NoError(t, err)

	err = billing.NextState()
	assert.NoError(t, err)

	err = billing.NextState()
	assert.NoError(t, err)

	err = billing.NextState()
	assert.EqualError(t, err, (ErrNextCompletedState{}).Error())

	state = billing.GetState()
	assert.Equal(t, StateCompleted, state)
}
