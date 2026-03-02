package commands

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/mocks/core/domain/servicesmocks"
	"basket-service/mocks/core/portsmocks"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_CheckoutCommandHandlerShouldBeSuccessWhenParamsAreCorrect(t *testing.T) {

	// Arrange
	ctx := context.Background()
	basketAggregate, err := basket.NewBasket(uuid.New())
	assert.NoError(t, err)
	address, err := kernel.NewAddress("Россия", "Москва", "Тверская", "1", "2")
	assert.NoError(t, err)
	deliveryPeriod := basket.Morning
	err = basketAggregate.Change(good.Bread(), 1)
	assert.NoError(t, err)
	err = basketAggregate.AddAddress(address)
	assert.NoError(t, err)
	err = basketAggregate.AddDeliveryPeriod(deliveryPeriod)
	assert.NoError(t, err)
	discount, err := kernel.NewDiscount(0.01)
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

	discountClientMock := &portsmocks.DiscountClientMock{}
	discountClientMock.
		On("GetDiscount", ctx, basketAggregate).
		Return(discount, nil).
		Once()

	promoGoodServiceMock := &servicesmocks.PromoGoodServiceMock{}
	promoGoodServiceMock.
		On("AddPromo", basketAggregate).
		Return(nil).
		Once()

	// Act
	checkoutCommandHandler, err := NewCheckoutCommandHandler(
		unitOfWorkMock, promoGoodServiceMock, discountClientMock)
	assert.NoError(t, err)
	checkoutCommand, err := NewCheckoutCommand(basketAggregate.ID())
	assert.NoError(t, err)
	err = checkoutCommandHandler.Handle(ctx, checkoutCommand)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedObj)
	assert.Equal(t, 99.0, capturedObj.Total())
	assert.Equal(t, basket.StatusConfirmed, capturedObj.Status())
}
