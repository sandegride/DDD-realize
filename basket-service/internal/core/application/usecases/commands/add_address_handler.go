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
	uowFactory ports.UnitOfWorkFactory
}

func NewAddAddressCommandHandler(
	uowFactory ports.UnitOfWorkFactory) (AddAddressCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &addAddressCommandHandler{uowFactory: uowFactory}, nil
}

func (ch *addAddressCommandHandler) Handle(ctx context.Context, command AddAddressCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add address command")
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
	address, err := kernel.NewAddress(command.Country(), command.City(), command.Street(),
		command.House(), command.Apartment())
	if err != nil {
		return err
	}
	err = basketAggregate.AddAddress(address)
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
