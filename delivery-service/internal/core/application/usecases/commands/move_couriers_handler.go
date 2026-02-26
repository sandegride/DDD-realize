package commands

import (
	"context"
	"delivery-service/internal/core/ports"
	"delivery-service/internal/pkg/errs"
	"errors"
)

type MoveCouriersCommandHandler interface {
	Handle(context.Context, MoveCouriersCommand) error
}

var _ MoveCouriersCommandHandler = &moveCouriersCommandHandler{}

type moveCouriersCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewMoveCouriersCommandHandler(
	uowFactory ports.UnitOfWorkFactory) (MoveCouriersCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &moveCouriersCommandHandler{
		uowFactory: uowFactory}, nil
}

func (ch *moveCouriersCommandHandler) Handle(ctx context.Context, command MoveCouriersCommand) error {
	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	assignedOrders, err := uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	uow.Begin(ctx)
	for _, assignedOrder := range assignedOrders {
		courier, err := uow.CourierRepository().Get(ctx, *assignedOrder.CourierID())
		if err != nil {
			if errors.Is(err, errs.ErrObjectNotFound) {
				return nil
			}
			return err
		}

		err = courier.Move(assignedOrder.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equals(assignedOrder.Location()) {
			err := assignedOrder.Complete()
			if err != nil {
				return err
			}
			err = courier.CompleteOrder(assignedOrder)
			if err != nil {
				return err
			}
		}

		err = uow.OrderRepository().Update(ctx, assignedOrder)
		if err != nil {
			return err
		}
		err = uow.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}
	}
	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
