package queries

import (
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type GetBasketQuery struct {
	basketID uuid.UUID

	isValid bool
}

func (q GetBasketQuery) BasketID() uuid.UUID {
	return q.basketID
}

func NewGetBasketQuery(basketID uuid.UUID) (GetBasketQuery, error) {
	if basketID == uuid.Nil {
		return GetBasketQuery{}, errs.NewValueIsInvalidError("basketID")
	}

	return GetBasketQuery{
		basketID: basketID,

		isValid: true,
	}, nil
}

func (q GetBasketQuery) IsValid() bool {
	return q.isValid
}
