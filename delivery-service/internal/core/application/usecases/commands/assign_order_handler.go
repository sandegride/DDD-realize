package commands

import (
	"context"
	"delivery-service/internal/core/domain/services"
	"delivery-service/internal/core/ports"
	"delivery-service/internal/pkg/errs"
	"errors"
)

var (
	NotAvailableOrders   = errors.New("not available orders")
	NotAvailableCouriers = errors.New("not available couriers")
)

type AssignOrdersCommandHandler interface {
	Handle(context.Context, AssignOrdersCommand) error
}

var _ AssignOrdersCommandHandler = &assignOrdersCommandHandler{}

type assignOrdersCommandHandler struct {
	uowFactory      ports.UnitOfWorkFactory
	orderDispatcher services.OrderDispatcherService
}

func NewAssignOrdersCommandHandler(
	uowFactory ports.UnitOfWorkFactory,
	orderDispatcher services.OrderDispatcherService) (AssignOrdersCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return &assignOrdersCommandHandler{
		uowFactory:      uowFactory,
		orderDispatcher: orderDispatcher}, nil
}

func (ch *assignOrdersCommandHandler) Handle(ctx context.Context, command AssignOrdersCommand) error {
	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	orderAggregate, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return NotAvailableOrders
		}
		return err
	}

	couriers, err := uow.CourierRepository().GetAllFree(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return NotAvailableCouriers
		}
		return err
	}
	if len(couriers) == 0 {
		return nil
	}

	courier, err := ch.orderDispatcher.Dispatch(orderAggregate, couriers)
	if err != nil {
		return err
	}

	uow.Begin(ctx)

	err = uow.OrderRepository().Update(ctx, orderAggregate)
	if err != nil {
		return err
	}
	err = uow.CourierRepository().Update(ctx, courier)
	if err != nil {
		return err
	}

	err = uow.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
