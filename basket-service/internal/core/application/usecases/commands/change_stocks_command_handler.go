package commands

import (
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
)

type ChangeStocksCommandHandler interface {
	Handle(context.Context, ChangeStocksCommand) error
}

var _ ChangeStocksCommandHandler = &changeStocksCommandHandler{}

type changeStocksCommandHandler struct {
	unitOfWork ports.UnitOfWork
}

func NewChangeStocksCommandHandler(unitOfWork ports.UnitOfWork) (ChangeStocksCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &changeStocksCommandHandler{
		unitOfWork: unitOfWork}, nil
}

func (ch *changeStocksCommandHandler) Handle(ctx context.Context, command ChangeStocksCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("change stocks command")
	}

	// Восстановили
	goodAggregate, err := ch.unitOfWork.GoodRepository().Get(ctx, command.GoodID())
	if err != nil {
		return err
	}

	// Изменили
	err = goodAggregate.ChangeStocks(command.Quantity())
	if err != nil {
		return err
	}

	// Сохранили
	err = ch.unitOfWork.GoodRepository().Update(ctx, goodAggregate)
	if err != nil {
		return err
	}
	return nil
}
