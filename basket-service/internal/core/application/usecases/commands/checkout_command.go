package commands

import (
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type CheckoutCommand struct {
	basketID uuid.UUID

	isValid bool
}

func (c CheckoutCommand) BasketID() uuid.UUID {
	return c.basketID

}

func NewCheckoutCommand(basketID uuid.UUID) (CheckoutCommand, error) {
	if basketID == uuid.Nil {
		return CheckoutCommand{}, errs.NewValueIsInvalidError("basketID")
	}

	return CheckoutCommand{
		basketID: basketID,
		isValid:  true,
	}, nil
}

func (c CheckoutCommand) IsValid() bool {
	return c.isValid
}
