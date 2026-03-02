package services

import (
	"delivery-service/internal/core/domain/model/courier"
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/domain/model/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DispatchService_ShouldFindNearestCourierForOrder(t *testing.T) {
	// Arrange
	courier1Location, err := kernel.NewLocation(1, 1)
	assert.NoError(t, err)
	courier1, err := courier.NewCourier("Pedestrian 1", 1, courier1Location)
	assert.NoError(t, err)
	err = courier1.AddStoragePlace("Сумка", 10)
	assert.NoError(t, err)

	courier2Location, err := kernel.NewLocation(2, 2)
	assert.NoError(t, err)
	courier2, err := courier.NewCourier("Pedestrian 2", 1, courier2Location)
	assert.NoError(t, err)
	err = courier2.AddStoragePlace("Сумка", 10)
	assert.NoError(t, err)

	courier3Location, err := kernel.NewLocation(3, 3)
	assert.NoError(t, err)
	courier3, err := courier.NewCourier("Pedestrian 3", 1, courier3Location)
	assert.NoError(t, err)
	err = courier3.AddStoragePlace("Сумка", 10)
	assert.NoError(t, err)

	orderLocation, err := kernel.NewLocation(2, 2)
	assert.NoError(t, err)
	orderAggregate, err := order.NewOrder(uuid.New(), orderLocation, 5)
	assert.NoError(t, err)

	couriers := []*courier.Courier{courier1, courier2, courier3}

	// Act
	dispatchService := NewOrderDispatcher()
	winner, err := dispatchService.Dispatch(orderAggregate, couriers)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, courier2, winner)
	assert.Equal(t, orderAggregate.ID(), *courier2.StoragePlaces()[0].OrderID())
	assert.Equal(t, order.StatusAssigned, orderAggregate.Status())
	assert.Equal(t, courier2.ID(), *orderAggregate.CourierID())
}

func Test_DispatchService_ShouldReturnError_WhenOrderIsNil(t *testing.T) {
	dispatcher := NewOrderDispatcher()
	_, err := dispatcher.Dispatch(nil, []*courier.Courier{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order")
}

func Test_DispatchService_ShouldReturnError_WhenCouriersListIsEmpty(t *testing.T) {
	location := kernel.MinLocation()
	o, err := order.NewOrder(uuid.New(), location, 1)
	assert.NoError(t, err)

	dispatcher := NewOrderDispatcher()
	_, err = dispatcher.Dispatch(o, []*courier.Courier{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "couriers")
}

func Test_DispatchService_ShouldReturnError_WhenNoCourierCanTakeOrder(t *testing.T) {
	location := kernel.MinLocation()
	o, err := order.NewOrder(uuid.New(), location, 20)
	assert.NoError(t, err)

	c1, err := courier.NewCourier("Пеший с сумкой 10 литров", 1, location)
	assert.NoError(t, err)

	dispatcher := NewOrderDispatcher()
	_, err = dispatcher.Dispatch(o, []*courier.Courier{c1})
	assert.Error(t, err)
	assert.Equal(t, ErrSuitableCourierWasNotFound, err)
}

func Test_DispatchService_ShouldSelectFirstCourier_WhenEqualDistance(t *testing.T) {
	orderLoc, err := kernel.NewLocation(2, 2)
	assert.NoError(t, err)
	o, err := order.NewOrder(uuid.New(), orderLoc, 1)
	assert.NoError(t, err)

	location := kernel.MinLocation()
	c1, err := courier.NewCourier("First", 1, location)
	assert.NoError(t, err)
	c2, err := courier.NewCourier("Second", 1, location)
	assert.NoError(t, err)

	dispatcher := NewOrderDispatcher()
	winner, err := dispatcher.Dispatch(o, []*courier.Courier{c1, c2})
	assert.NoError(t, err)
	assert.Equal(t, c1, winner)
}
