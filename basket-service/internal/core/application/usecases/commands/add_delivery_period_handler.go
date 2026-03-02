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
	unitOfWork ports.UnitOfWork
}

func NewAddDeliveryPeriodCommandHandler(unitOfWork ports.UnitOfWork) (AddDeliveryPeriodCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &addDeliveryPeriodCommandHandler{
		unitOfWork: unitOfWork}, nil
}

func (ch *addDeliveryPeriodCommandHandler) Handle(ctx context.Context, command AddDeliveryPeriodCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add delivery period command")
	}

	basketAggregate, err := ch.unitOfWork.BasketRepository().Get(ctx, command.BasketID())
	if err != nil {
		return err
	}

	deliveryPeriod, err := basket.GetDeliveryPeriodByName(command.DeliveryPeriod().String())
	if err != nil {
		return err
	}
	err = basketAggregate.AddDeliveryPeriod(deliveryPeriod)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.BasketRepository().Update(ctx, basketAggregate)
	if err != nil {
		return err
	}

	return nil
}
