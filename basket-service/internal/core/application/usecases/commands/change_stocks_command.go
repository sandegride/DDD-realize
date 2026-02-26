package commands

import (
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type ChangeStocksCommand struct {
	goodID   uuid.UUID
	quantity int

	isValid bool
}

func (c ChangeStocksCommand) GoodID() uuid.UUID {
	return c.goodID
}

func (c ChangeStocksCommand) Quantity() int {
	return c.quantity
}

func NewChangeStocksCommand(goodID uuid.UUID, quantity int) (ChangeStocksCommand, error) {
	if goodID == uuid.Nil {
		return ChangeStocksCommand{}, errs.NewValueIsRequiredError("goodID")
	}
	if quantity < 0 {
		return ChangeStocksCommand{}, errs.NewValueIsInvalidError("quantity")
	}

	return ChangeStocksCommand{
		goodID:   goodID,
		quantity: quantity,
		isValid:  true,
	}, nil
}

func (c ChangeStocksCommand) IsValid() bool {
	return c.isValid
}
