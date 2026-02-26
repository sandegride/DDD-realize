package commands

import (
	"context"
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/domain/model/order"
	portsmocks "delivery-service/internal/mocks/core/portmocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_CreateOrderHandlerShouldBeSuccessWhenOrderDoesNotExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orderID := uuid.New()
	street := "street"
	volume := 12

	cmd, err := NewCreateOrderCommand(orderID, street, volume)
	assert.NoError(t, err)

	orderRepoMock := &portsmocks.OrderRepositoryMock{}
	orderRepoMock.
		On("Get", ctx, orderID).
		Return(nil, nil).
		Once()

	var capturedOrder *order.Order
	orderRepoMock.
		On("Add", ctx, mock.AnythingOfType("*order.Order")).
		Run(func(args mock.Arguments) {
			capturedOrder = args.Get(1).(*order.Order)
		}).
		Return(nil).
		Once()

	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.
		On("OrderRepository").
		Return(orderRepoMock)
	uowMock.
		On("RollbackUnlessCommitted", ctx).
		Return()

	uowFactoryMock := &portsmocks.UnitOfWorkFactoryMock{}
	uowFactoryMock.
		On("New", ctx).
		Return(uowMock, nil)

	handler, err := NewCreateOrderCommandHandler(uowFactoryMock)
	assert.NoError(t, err)

	// Act
	err = handler.Handle(ctx, cmd)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedOrder, "заказ должен быть передан в репозиторий")
	assert.Equal(t, orderID, capturedOrder.ID())
	assert.Equal(t, volume, capturedOrder.Volume())
	assert.NotNil(t, capturedOrder.Location(), "локация должна быть сгенерирована")
}

func Test_CreateOrderHandlerShouldDoNothingWhenOrderAlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orderID := uuid.New()
	street := "street"
	volume := 7

	existingOrder, err := order.NewOrder(orderID, kernel.NewRandomLocation(), volume)
	assert.NoError(t, err)

	cmd, err := NewCreateOrderCommand(orderID, street, volume)
	assert.NoError(t, err)

	orderRepoMock := &portsmocks.OrderRepositoryMock{}
	orderRepoMock.
		On("Get", ctx, orderID).
		Return(existingOrder, nil).
		Once()

	uowMock := &portsmocks.UnitOfWorkMock{}
	uowMock.
		On("OrderRepository").
		Return(orderRepoMock)
	uowMock.
		On("RollbackUnlessCommitted", ctx).
		Return()

	uowFactoryMock := &portsmocks.UnitOfWorkFactoryMock{}
	uowFactoryMock.
		On("New", ctx).
		Return(uowMock, nil)

	handler, err := NewCreateOrderCommandHandler(uowFactoryMock)
	assert.NoError(t, err)

	// Act
	err = handler.Handle(ctx, cmd)

	// Assert
	assert.NoError(t, err)
	orderRepoMock.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
}
