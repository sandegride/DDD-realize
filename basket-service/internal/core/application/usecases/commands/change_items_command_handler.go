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
	unitOfWork ports.UnitOfWork
}

func NewChangeItemsCommandHandler(unitOfWork ports.UnitOfWork) (ChangeItemsCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &changeItemsCommandHandler{
		unitOfWork: unitOfWork}, nil
}

func (ch *changeItemsCommandHandler) Handle(ctx context.Context, command ChangeItemsCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("change items command")
	}

	// Восстановили
	goodAggregate, err := ch.unitOfWork.GoodRepository().Get(ctx, command.GoodID())
	if err != nil {
		return err
	}

	basketAggregate, err := ch.unitOfWork.BasketRepository().Get(ctx, command.BasketID())
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			basketAggregate, err = basket.NewBasket(command.BasketID())
			if err != nil {
				return err
			}

			// Добавили
			err = ch.unitOfWork.BasketRepository().Add(ctx, basketAggregate)
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
	ch.unitOfWork.Begin(ctx)

	// Обновили данные в репозиториях
	err = ch.unitOfWork.BasketRepository().Update(ctx, basketAggregate)
	if err != nil {
		return err
	}
	err = ch.unitOfWork.GoodRepository().Update(ctx, goodAggregate)
	if err != nil {
		return err
	}

	// Завершили транзакцию
	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
