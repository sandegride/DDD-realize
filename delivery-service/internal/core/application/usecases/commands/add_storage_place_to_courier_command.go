package commands

import (
	"delivery-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type AddStoragePlaceToCourierCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int

	isValid bool
}

func NewAddStoragePlaceCourierCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlaceToCourierCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceToCourierCommand{}, errs.NewValueIsRequiredError("courierID")
	}
	if name == "" {
		return AddStoragePlaceToCourierCommand{}, errs.NewValueIsRequiredError("name")
	}
	if totalVolume <= 0 {
		return AddStoragePlaceToCourierCommand{}, errs.NewValueIsRequiredError("totalVolume")
	}

	return AddStoragePlaceToCourierCommand{courierID: courierID, name: name, totalVolume: totalVolume, isValid: true}, nil
}

func (c AddStoragePlaceToCourierCommand) CourierID() uuid.UUID {
	return c.courierID
}

func (c AddStoragePlaceToCourierCommand) Name() string {
	return c.name
}

func (c AddStoragePlaceToCourierCommand) TotalVolume() int {
	return c.totalVolume
}

func (c AddStoragePlaceToCourierCommand) IsValid() bool {
	return c.isValid
}
