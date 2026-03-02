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
	unitOfWork      ports.UnitOfWork
	orderDispatcher services.OrderDispatcher
}

func NewAssignOrdersCommandHandler(
	unitOfWork ports.UnitOfWork,
	orderDispatcher services.OrderDispatcher) (AssignOrdersCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return &assignOrdersCommandHandler{
		unitOfWork:      unitOfWork,
		orderDispatcher: orderDispatcher}, nil
}

func (ch *assignOrdersCommandHandler) Handle(ctx context.Context, command AssignOrdersCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add address command")
	}

	orderAggregate, err := ch.unitOfWork.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return NotAvailableOrders
		}
		return err
	}

	couriers, err := ch.unitOfWork.CourierRepository().GetAllFree(ctx)
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

	ch.unitOfWork.Begin(ctx)

	err = ch.unitOfWork.OrderRepository().Update(ctx, orderAggregate)
	if err != nil {
		return err
	}
	err = ch.unitOfWork.CourierRepository().Update(ctx, courier)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
