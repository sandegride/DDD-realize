package commands

import (
	"context"
	"delivery-service/internal/core/ports"
	"delivery-service/internal/pkg/errs"
)

type AddStoragePlaceToCourierCommandHandler interface {
	Handle(ctx context.Context, command AddStoragePlaceToCourierCommand) error
}

var _ AddStoragePlaceToCourierCommandHandler = &addStoragePlaceToCourierCommandHandler{}

type addStoragePlaceToCourierCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewAddStoragePlaceToCourierCommandHandler(uowFactory ports.UnitOfWorkFactory) (AddStoragePlaceToCourierCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &addStoragePlaceToCourierCommandHandler{uowFactory: uowFactory}, nil
}

func (ch addStoragePlaceToCourierCommandHandler) Handle(ctx context.Context, command AddStoragePlaceToCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add storage place to courier command")
	}

	uof, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uof.RollbackUnlessCommitted(ctx)

	courierAggregate, err := uof.CourierRepository().Get(ctx, command.CourierID())
	if err != nil {
		return err
	}

	err = courierAggregate.AddStoragePlace(command.Name(), command.TotalVolume())
	if err != nil {
		return err
	}

	err = uof.CourierRepository().Update(ctx, courierAggregate)
	if err != nil {
		return err
	}

	return nil
}
