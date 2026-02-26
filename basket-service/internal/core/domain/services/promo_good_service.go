package services

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/pkg/errs"
)

type PromoGoodService interface {
	AddPromo(basket *basket.Basket) error
}

var _ PromoGoodService = &promoGoodService{}

type promoGoodService struct{}

func NewPromoGoodService() PromoGoodService {
	return &promoGoodService{}
}

func (p *promoGoodService) AddPromo(basket *basket.Basket) error {
	if basket == nil {
		return errs.NewValueIsRequiredError("basket")
	}

	promoGum := good.Gum()
	promoCandy := good.Candy()
	promoSnack := good.Snack()
	total := basket.Total()

	switch {
	case total > 1000 && total <= 2000:
		err := basket.Change(promoGum, 1)
		if err != nil {
			return err
		}
		break
	case total > 2000 && total <= 5000:
		err := basket.Change(promoCandy, 1)
		if err != nil {
			return err
		}
		break
	case total > 5000:
		err := basket.Change(promoSnack, 1)
		if err != nil {
			return err
		}
		break
	}
	return nil
}
