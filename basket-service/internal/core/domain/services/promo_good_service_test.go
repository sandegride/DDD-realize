package services

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/good"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PromoGoodServiceShouldAddGumToBasketWhenTotalIsMoreThen1000(t *testing.T) {
	buyerID := uuid.New()
	basketAggregate, _ := basket.NewBasket(buyerID)
	_ = basketAggregate.Change(good.Bread(), 11) // 100 * 11= 1100

	promoGoodService := NewPromoGoodService()
	err := promoGoodService.AddPromo(basketAggregate) // gum for 1 ruble

	assert.NoError(t, err)
	assert.Equal(t, good.Bread().ID(), basketAggregate.Items()[0].GoodID())
	assert.Equal(t, good.Gum().ID(), basketAggregate.Items()[1].GoodID())
	assert.Equal(t, 1101.0, basketAggregate.Total())
}
