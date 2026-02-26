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
	uowFactory ports.UnitOfWorkFactory
}

func NewChangeStocksCommandHandler(
	uowFactory ports.UnitOfWorkFactory) (ChangeStocksCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &changeStocksCommandHandler{
		uowFactory: uowFactory}, nil
}

func (ch *changeStocksCommandHandler) Handle(ctx context.Context, command ChangeStocksCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("change stocks command")
	}

	// Создаем UoW
	uow, err := ch.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	// Восстановили
	goodAggregate, err := uow.GoodRepository().Get(ctx, command.GoodID())
	if err != nil {
		return err
	}

	// Изменили
	err = goodAggregate.ChangeStocks(command.Quantity())
	if err != nil {
		return err
	}

	// Сохранили
	err = uow.GoodRepository().Update(ctx, goodAggregate)
	if err != nil {
		return err
	}
	return nil
}
