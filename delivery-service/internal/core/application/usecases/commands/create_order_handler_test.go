package commands

import (
	"context"
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/domain/model/order"
	"delivery-service/mocks/core/portsmocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_CreateOrderCommandHandlerShouldBeSuccessWhenParamsAreCorrect(t *testing.T) {
	// Arrange
	ctx := context.Background()
	location := kernel.MinLocation()
	orderAggregate, err := order.NewOrder(uuid.New(), location, 10)
	assert.NoError(t, err)

	var capturedObj *order.Order
	orderRepositoryMock := &portsmocks.OrderRepositoryMock{}
	orderRepositoryMock.
		On("Get", ctx, orderAggregate.ID()).
		Return(nil, nil).
		Once()
	orderRepositoryMock.
		On("Add", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			capturedObj = args.Get(1).(*order.Order)
		}).
		Return(nil, nil).
		Once()

	unitOfWorkMock := &portsmocks.UnitOfWorkMock{}
	unitOfWorkMock.
		On("OrderRepository").
		Return(orderRepositoryMock)
	unitOfWorkFactoryMock := &portsmocks.UnitOfWorkFactoryMock{}
	unitOfWorkFactoryMock.
		On("New", ctx).
		Return(unitOfWorkMock, nil)
	unitOfWorkMock.
		On("RollbackUnlessCommitted", ctx).
		Return()

	geoClientMock := &portsmocks.GeoClientMock{}
	geoClientMock.
		On("GetGeolocation", ctx, mock.Anything).
		Return(location, nil).
		Once()

	// Act
	createOrderCommandHandler, err := NewCreateOrderCommandHandler(unitOfWorkFactoryMock, geoClientMock)
	assert.NoError(t, err)
	createOrderCommand, err := NewCreateOrderCommand(orderAggregate.ID(),
		"Несуществующая", 10)
	assert.NoError(t, err)
	err = createOrderCommandHandler.Handle(ctx, createOrderCommand)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedObj.ID())
	assert.True(t, capturedObj.Location().IsValid())
	assert.Equal(t, order.StatusCreated, capturedObj.Status())
}
