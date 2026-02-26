package order

import (
	"delivery-service/internal/core/domain/model/kernel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrderShouldBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	orderID := uuid.New()
	location := kernel.MinLocation()

	order, err := NewOrder(orderID, location, 10)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orderID, order.ID())
	assert.Equal(t, uuid.Nil, order.CourierID())
	assert.Equal(t, location, order.Location())
	assert.Equal(t, 10, order.Volume())
	assert.Equal(t, StatusCreated, order.Status())
}
