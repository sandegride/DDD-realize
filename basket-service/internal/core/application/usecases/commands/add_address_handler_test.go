package commands

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/mocks/core/portsmocks"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_AddAddressCommandHandlerShouldBeSuccessWhenParamsAreCorrect(t *testing.T) {
	// Arrange
	ctx := context.Background()
	basketAggregate, err := basket.NewBasket(uuid.New())
	assert.NoError(t, err)

	var capturedObj *basket.Basket
	basketRepositoryMock := &portsmocks.BasketRepositoryMock{}
	basketRepositoryMock.
		On("Get", ctx, basketAggregate.ID()).
		Return(basketAggregate, nil).
		Once()
	basketRepositoryMock.
		On("Update", ctx, basketAggregate).
		Run(func(args mock.Arguments) {
			capturedObj = args.Get(1).(*basket.Basket)
		}).
		Return(nil, nil).
		Once()
	unitOfWorkMock := &portsmocks.UnitOfWorkMock{}
	unitOfWorkMock.
		On("BasketRepository").
		Return(basketRepositoryMock)

	// Act
	addAddressCommandHandler, err := NewAddAddressCommandHandler(unitOfWorkMock)
	assert.NoError(t, err)
	addAddressCommand, err := NewAddAddressCommand(basketAggregate.ID(),
		"Россия", "Москва", "Тверская", "1", "2")
	assert.NoError(t, err)
	err = addAddressCommandHandler.Handle(ctx, addAddressCommand)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Россия", capturedObj.Address().Country())
	assert.Equal(t, "Москва", capturedObj.Address().City())
	assert.Equal(t, "Тверская", capturedObj.Address().Street())
	assert.Equal(t, "1", capturedObj.Address().House())
	assert.Equal(t, "2", capturedObj.Address().Apartment())
}
