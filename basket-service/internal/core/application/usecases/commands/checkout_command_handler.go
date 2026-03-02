package commands

import (
	"basket-service/internal/core/domain/services"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
)

type CheckoutCommandHandler interface {
	Handle(context.Context, CheckoutCommand) error
}

var _ CheckoutCommandHandler = &checkoutCommandHandler{}

type checkoutCommandHandler struct {
	unitOfWork       ports.UnitOfWork
	promoGoodService services.PromoGoodService
	discountClient   ports.DiscountClient
}

func NewCheckoutCommandHandler(
	unitOfWork ports.UnitOfWork,
	promoGoodService services.PromoGoodService, discountClient ports.DiscountClient) (CheckoutCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if promoGoodService == nil {
		return nil, errs.NewValueIsRequiredError("promoGoodService")
	}
	if discountClient == nil {
		return nil, errs.NewValueIsRequiredError("discountClient")
	}

	return &checkoutCommandHandler{
		unitOfWork:       unitOfWork,
		promoGoodService: promoGoodService,
		discountClient:   discountClient,
	}, nil
}

func (ch *checkoutCommandHandler) Handle(ctx context.Context, command CheckoutCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("checkout command")
	}

	basketAggregate, err := ch.unitOfWork.BasketRepository().Get(ctx, command.BasketID())
	if err != nil {
		return err
	}

	discount, err := ch.discountClient.GetDiscount(ctx, basketAggregate)
	if err != nil {
		return err
	}

	err = ch.promoGoodService.AddPromo(basketAggregate)
	if err != nil {
		return err
	}

	err = basketAggregate.Checkout(discount)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.BasketRepository().Update(ctx, basketAggregate)
	if err != nil {
		return err
	}

	return nil
}
