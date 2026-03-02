package commands

import (
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
)

type AddAddressCommandHandler interface {
	Handle(context.Context, AddAddressCommand) error
}

var _ AddAddressCommandHandler = &addAddressCommandHandler{}

type addAddressCommandHandler struct {
	unitOfWork ports.UnitOfWork
}

func NewAddAddressCommandHandler(
	unitOfWork ports.UnitOfWork) (AddAddressCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &addAddressCommandHandler{
		unitOfWork: unitOfWork}, nil
}

func (ch *addAddressCommandHandler) Handle(ctx context.Context, command AddAddressCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add address command")
	}

	basketAggregate, err := ch.unitOfWork.BasketRepository().Get(ctx, command.BasketID())
	if err != nil {
		return err
	}

	address, err := kernel.NewAddress(command.Country(), command.City(), command.Street(),
		command.House(), command.Apartment())
	if err != nil {
		return err
	}
	err = basketAggregate.AddAddress(address)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.BasketRepository().Update(ctx, basketAggregate)
	if err != nil {
		return err
	}

	return nil
}
