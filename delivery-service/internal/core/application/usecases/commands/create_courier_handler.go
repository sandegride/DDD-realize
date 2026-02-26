package commands

import (
	"context"
	"delivery-service/internal/core/domain/model/courier"
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/ports"
	"delivery-service/internal/pkg/errs"
)

type CreateCourierHandler interface {
	Handle(ctx context.Context, command CreateCourierCommand) error
}

var _ CreateCourierHandler = &createCourierHandler{}

type createCourierHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateCourierHandler(uowFactory ports.UnitOfWorkFactory) (CreateCourierHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &createCourierHandler{uowFactory: uowFactory}, nil
}

func (h *createCourierHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("create courier command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	location := kernel.NewRandomLocation()
	courierAggregate, err := courier.NewCourier(command.Name(), command.Speed(), location)
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Add(ctx, courierAggregate)
	if err != nil {
		return err
	}
	return nil

}
