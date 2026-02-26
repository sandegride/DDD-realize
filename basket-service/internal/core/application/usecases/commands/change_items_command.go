package commands

import (
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type ChangeItemsCommand struct {
	basketID uuid.UUID
	buyerID  uuid.UUID
	goodID   uuid.UUID
	quantity int

	isValid bool
}

func (c ChangeItemsCommand) BasketID() uuid.UUID {
	return c.basketID
}

func (c ChangeItemsCommand) BuyerID() uuid.UUID {
	return c.buyerID
}

func (c ChangeItemsCommand) GoodID() uuid.UUID {
	return c.goodID
}

func (c ChangeItemsCommand) Quantity() int {
	return c.quantity
}

func NewChangeItemsCommand(basketID uuid.UUID, buyerID uuid.UUID, goodID uuid.UUID, quantity int) (ChangeItemsCommand, error) {
	if basketID == uuid.Nil {
		return ChangeItemsCommand{}, errs.NewValueIsInvalidError("basketID")
	}
	if buyerID == uuid.Nil {
		return ChangeItemsCommand{}, errs.NewValueIsInvalidError("buyerID")
	}
	if goodID == uuid.Nil {
		return ChangeItemsCommand{}, errs.NewValueIsInvalidError("goodID")
	}
	if quantity <= 0 {
		return ChangeItemsCommand{}, errs.NewValueIsInvalidError("quantity")
	}

	return ChangeItemsCommand{
		basketID: basketID,
		buyerID:  buyerID,
		goodID:   goodID,
		quantity: quantity,
		isValid:  true,
	}, nil
}

func (c ChangeItemsCommand) IsValid() bool {
	return c.isValid
}
