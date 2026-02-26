package commands

import (
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type AddDeliveryPeriodCommand struct {
	basketID       uuid.UUID
	deliveryPeriod DeliveryPeriod

	isValid bool
}

func (a AddDeliveryPeriodCommand) BasketID() uuid.UUID {
	return a.basketID
}

func (a AddDeliveryPeriodCommand) DeliveryPeriod() DeliveryPeriod {
	return a.deliveryPeriod
}

func NewAddDeliveryPeriodCommand(basketID uuid.UUID, deliveryPeriod DeliveryPeriod) (AddDeliveryPeriodCommand, error) {
	if basketID == uuid.Nil {
		return AddDeliveryPeriodCommand{}, errs.NewValueIsInvalidError("basketID")
	}
	if deliveryPeriod == "" {
		return AddDeliveryPeriodCommand{}, errs.NewValueIsRequiredError("deliveryPeriod")
	}

	return AddDeliveryPeriodCommand{
		basketID:       basketID,
		deliveryPeriod: deliveryPeriod,
		isValid:        true,
	}, nil
}

func (a AddDeliveryPeriodCommand) IsValid() bool {
	return a.isValid
}

const (
	DeliveryPeriodEmpty DeliveryPeriod = ""
)

type DeliveryPeriod string

func (d DeliveryPeriod) Equal(other DeliveryPeriod) bool {
	return d == other
}

func (d DeliveryPeriod) IsValid() bool {
	return d == DeliveryPeriodEmpty
}

func (d DeliveryPeriod) String() string {
	return string(d)
}
