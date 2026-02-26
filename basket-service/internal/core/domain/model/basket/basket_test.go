package basket

import (
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BasketShouldBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange
	buyerID := uuid.New()

	// Act
	basket, err := NewBasket(buyerID)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, basket.ID)
	assert.Equal(t, buyerID, basket.BuyerID())
	assert.Empty(t, basket.Address())
	assert.Nil(t, basket.DeliveryPeriod())
	assert.Equal(t, StatusCreated, basket.Status())
}

func Test_BasketShouldReturnValueIsRequiredErrorWhenBuyerIdIsEmpty(t *testing.T) {
	// Arrange
	buyerID := uuid.Nil

	// Act
	basket, err := NewBasket(buyerID)

	// Assert
	assert.Nil(t, basket)
	assert.Error(t, err)
	if err.Error() != errs.NewValueIsRequiredError("buyerID").Error() {
		t.Errorf("expected %v, got %v", errs.NewValueIsRequiredError("buyerID"), err)
	}
}

func Test_BasketShouldHave3ItemWhenGoodsIsDifferent(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)

	// Act
	err := basket.Change(good.Coffee(), 1)
	assert.NoError(t, err)
	err = basket.Change(good.Milk(), 2)
	assert.NoError(t, err)
	err = basket.Change(good.Sugar(), 3)
	assert.NoError(t, err)

	// Assert
	coffeeItem, _, err := basket.FindItemByGoodID(good.Coffee().ID())
	assert.NoError(t, err)
	milkItem, _, err := basket.FindItemByGoodID(good.Milk().ID())
	assert.NoError(t, err)
	sugarItem, _, err := basket.FindItemByGoodID(good.Sugar().ID())
	assert.NoError(t, err)

	assert.Equal(t, 3, len(basket.items))
	assert.Equal(t, 1, coffeeItem.quantity)
	assert.Equal(t, 2, milkItem.quantity)
	assert.Equal(t, 3, sugarItem.quantity)
}

func Test_BasketShouldHave2ItemWhenOneGoodsIsSame(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)

	// Act
	err := basket.Change(good.Coffee(), 1)
	assert.NoError(t, err)
	err = basket.Change(good.Coffee(), 2)
	assert.NoError(t, err)
	err = basket.Change(good.Sugar(), 3)
	assert.NoError(t, err)

	// Assert
	coffeeItem, _, err := basket.FindItemByGoodID(good.Coffee().ID())
	assert.NoError(t, err)
	sugarItem, _, err := basket.FindItemByGoodID(good.Sugar().ID())
	assert.NoError(t, err)

	assert.Equal(t, 2, len(basket.items))
	assert.Equal(t, 2, coffeeItem.quantity)
	assert.Equal(t, 3, sugarItem.quantity)
}

func Test_BasketShouldHave2ItemWhenOneGoodsIsZero(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)
	_ = basket.Change(good.Coffee(), 1)
	_ = basket.Change(good.Milk(), 2)
	_ = basket.Change(good.Sugar(), 3)

	// Act
	err := basket.Change(good.Coffee(), 0)
	assert.NoError(t, err)
	err = basket.Change(good.Milk(), 2)
	assert.NoError(t, err)
	err = basket.Change(good.Sugar(), 0)
	assert.NoError(t, err)

	// Assert
	coffeeItem, _, err := basket.FindItemByGoodID(good.Coffee().ID())
	assert.NoError(t, err)
	milkItem, _, err := basket.FindItemByGoodID(good.Milk().ID())
	assert.NoError(t, err)
	sugarItem, _, err := basket.FindItemByGoodID(good.Sugar().ID())
	assert.NoError(t, err)

	assert.Equal(t, 1, len(basket.items))
	assert.Nil(t, coffeeItem)
	assert.Equal(t, 2, milkItem.quantity)
	assert.Nil(t, sugarItem)
}

func Test_BasketShouldReturnValueIsRequiredErrorWhenAddNewItemWithZeroQuantity(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)

	// Act
	err := basket.Change(good.Coffee(), 0)

	// Assert
	assert.Equal(t, 0, len(basket.items))
	assert.Error(t, err)
	if err.Error() != ErrQuantityIsZeroOrLess.Error() {
		t.Errorf("expected %v, got %v", ErrQuantityIsZeroOrLess, err)
	}
}

func Test_BasketShouldCanAddAddress(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)
	_ = basket.Change(good.Coffee(), 1)
	_ = basket.Change(good.Milk(), 2)
	_ = basket.Change(good.Sugar(), 3)
	address, _ := kernel.NewAddress("Россия", "Москва", "Тверская", "1", "2")

	// Act
	err := basket.AddAddress(address)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, address, basket.Address())
}

func Test_BasketShouldCanAddDeliveryPeriod(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)
	_ = basket.Change(good.Coffee(), 1)
	_ = basket.Change(good.Milk(), 2)
	_ = basket.Change(good.Sugar(), 3)
	deliveryPeriod := Morning

	// Act
	err := basket.AddDeliveryPeriod(deliveryPeriod)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, Morning, basket.DeliveryPeriod())
}

func Test_BasketShouldCanCheckoutWhenBasketHasItemsAndDeliveryData(t *testing.T) {
	// Arrange
	buyerID := uuid.New()
	basket, _ := NewBasket(buyerID)
	_ = basket.Change(good.Coffee(), 1)
	_ = basket.Change(good.Milk(), 2)
	_ = basket.Change(good.Sugar(), 3)

	address, _ := kernel.NewAddress("Россия", "Москва", "Тверская", "1", "2")
	deliveryPeriod := Morning
	discount, _ := kernel.NewDiscount(0)

	_ = basket.AddAddress(address)
	_ = basket.AddDeliveryPeriod(deliveryPeriod)

	// Act
	err := basket.Checkout(discount)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StatusConfirmed, basket.Status())
}

func Test_BasketShouldHasCorrectTotalWithDiscount(t *testing.T) {
	// Arrange
	buyerID := uuid.New()

	address, _ := kernel.NewAddress("Россия", "Москва", "Тверская", "1", "2")
	deliveryPeriod := Morning
	zeroPercent, _ := kernel.NewDiscount(0.00)
	fivePercent, _ := kernel.NewDiscount(0.05)
	tenPercent, _ := kernel.NewDiscount(0.1)

	tests := map[string]struct {
		discount      kernel.Discount
		expectedTotal float64
	}{
		"0% -> 500": {
			discount:      zeroPercent,
			expectedTotal: 500,
		},
		"5% -> 475": {
			discount:      fivePercent,
			expectedTotal: 475,
		},
		"10% -> 450": {
			discount:      tenPercent,
			expectedTotal: 450,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			basket, _ := NewBasket(buyerID)
			_ = basket.Change(good.Bread(), 1)
			_ = basket.Change(good.Milk(), 2)
			_ = basket.AddAddress(address)
			_ = basket.AddDeliveryPeriod(deliveryPeriod)
			err := basket.Checkout(test.discount)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTotal, basket.Total())
		})
	}
}
