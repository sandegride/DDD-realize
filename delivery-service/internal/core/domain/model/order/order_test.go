package order

import (
	"delivery-service/internal/core/domain/model/kernel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrderShouldBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	location := kernel.MinLocation()

	// Act
	order, err := NewOrder(orderID, location, 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orderID, order.ID())
	assert.Equal(t, (*uuid.UUID)(nil), order.CourierID())
	assert.Equal(t, location, order.Location())
	assert.Equal(t, 10, order.Volume())
	assert.Equal(t, StatusCreated, order.Status())
}
