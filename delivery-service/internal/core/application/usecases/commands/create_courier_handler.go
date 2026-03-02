package commands

import (
	"context"
	"delivery-service/internal/core/domain/model/courier"
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/ports"
	"delivery-service/internal/pkg/errs"
)

type CreateCourierCommandHandler interface {
	Handle(context.Context, CreateCourierCommand) error
}

var _ CreateCourierCommandHandler = &createCourierCommandHandler{}

type createCourierCommandHandler struct {
	unitOfWork ports.UnitOfWork
}

func NewCreateCourierCommandHandler(
	unitOfWork ports.UnitOfWork) (CreateCourierCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &createCourierCommandHandler{
		unitOfWork: unitOfWork,
	}, nil
}

func (ch *createCourierCommandHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add address command")
	}

	location := kernel.NewRandomLocation()
	courierAggregate, err := courier.NewCourier(command.Name(), command.Speed(), location)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.CourierRepository().Add(ctx, courierAggregate)
	if err != nil {
		return err
	}
	return nil
}
