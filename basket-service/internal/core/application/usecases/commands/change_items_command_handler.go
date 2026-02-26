package commands

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
	"errors"
)

type ChangeItemsCommandHandler interface {
	Handle(context.Context, ChangeItemsCommand) error
}

var _ ChangeItemsCommandHandler = &changeItemsCommandHandler{}

type changeItemsCommandHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewChangeItemsCommandHandler(
	uowFactory ports.UnitOfWorkFactory) (ChangeItemsCommandHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &changeItemsCommandHandler{
		uowFactory: uowFactory}, nil
}

func (ch *changeItemsCommandHandler) Handle(ctx context.Context, command ChangeItemsCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("change items command")
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

	basketAggregate, err := uow.BasketRepository().Get(ctx, command.BasketID())
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			basketAggregate, err = basket.NewBasket(command.BasketID())
			if err != nil {
				return err
			}

			// Добавили
			err = uow.BasketRepository().Add(ctx, basketAggregate)
			if err != nil {
				return err
			}
		}
	}

	// Изменили
	err = basketAggregate.Change(goodAggregate, command.quantity)
	if err != nil {
		return err
	}

	// Изменили
	quantity := goodAggregate.Quantity() - command.quantity
	err = goodAggregate.ChangeStocks(quantity)
	if err != nil {
		return err
	}

	// Начали транзакцию
	uow.Begin(ctx)

	// Обновили данные в репозиториях
	err = uow.BasketRepository().Update(ctx, basketAggregate)
	if err != nil {
		return err
	}
	err = uow.GoodRepository().Update(ctx, goodAggregate)
	if err != nil {
		return err
	}

	// Завершили транзакцию
	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
