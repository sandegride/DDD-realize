package commands

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
)

type AddDeliveryPeriodCommandHandler interface {
	Handle(context.Context, AddDeliveryPeriodCommand) error
}

var _ AddDeliveryPeriodCommandHandler = &addDeliveryPeriodCommandHandler{}

type addDeliveryPeriodCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewAddDeliveryPeriodCommandHandler(
	uowFactory ports.UnitOfWorkFactory) (AddDeliveryPeriodCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &addDeliveryPeriodCommandHandler{
		uowFactory: uowFactory}, nil
}

func (ch *addDeliveryPeriodCommandHandler) Handle(ctx context.Context, command AddDeliveryPeriodCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add delivery period command")
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

	// Изменили
	deliveryPeriod, err := basket.GetDeliveryPeriodByName(command.DeliveryPeriod().String())
	if err != nil {
		return err
	}
	err = basketAggregate.AddDeliveryPeriod(deliveryPeriod)
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
