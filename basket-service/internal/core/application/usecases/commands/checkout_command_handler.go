package commands

import (
	"basket-service/internal/core/domain/model/kernel"
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
	uowFactory       ports.UnitOfWorkFactory
	promoGoodService services.PromoGoodService
}

func NewCheckoutCommandHandler(
	uowFactory ports.UnitOfWorkFactory,
	promoGoodService services.PromoGoodService) (CheckoutCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if promoGoodService == nil {
		return nil, errs.NewValueIsRequiredError("promoGoodService")
	}

	return &checkoutCommandHandler{
		uowFactory:       uowFactory,
		promoGoodService: promoGoodService,
	}, nil
}

func (ch *checkoutCommandHandler) Handle(ctx context.Context, command CheckoutCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("checkout command")
	}

	// Создаем UoW
	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	// Восстановили
	basketAggregate, err := uow.BasketRepository().Get(ctx, command.BasketID())
	if err != nil {
		return err
	}

	// Получаем геопозицию из сервиса Discount. Пока не реализовано - ставим дефолтное значение
	discount, err := kernel.NewDiscount(0.01)

	// Добавили промо товары (вызвали Domain Service)
	err = ch.promoGoodService.AddPromo(basketAggregate)
	if err != nil {
		return err
	}

	// Оформили корзину (изменили состояние Aggregate)
	err = basketAggregate.Checkout(discount)
	if err != nil {
		return err
	}

	// Сохранили
	err = uow.BasketRepository().Update(ctx, basketAggregate)
	if err != nil {
		return err
	}

	return nil
}
