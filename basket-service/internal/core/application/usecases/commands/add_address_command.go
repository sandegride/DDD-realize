package commands

import (
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

type AddAddressCommand struct {
	basketID  uuid.UUID
	country   string
	city      string
	street    string
	house     string
	apartment string

	isValid bool
}

func (c AddAddressCommand) BasketID() uuid.UUID {
	return c.basketID
}

func (c AddAddressCommand) Country() string {
	return c.country
}

func (c AddAddressCommand) City() string {
	return c.city
}

func (c AddAddressCommand) Street() string {
	return c.street
}

func (c AddAddressCommand) House() string {
	return c.house
}

func (c AddAddressCommand) Apartment() string {
	return c.apartment
}

func NewAddAddressCommand(basketID uuid.UUID, country string, city string, street string, house string, apartment string) (AddAddressCommand, error) {
	if basketID == uuid.Nil {
		return AddAddressCommand{}, errs.NewValueIsInvalidError("basketID")
	}
	if country == "" {
		return AddAddressCommand{}, errs.NewValueIsRequiredError("country")
	}
	if city == "" {
		return AddAddressCommand{}, errs.NewValueIsRequiredError("city")
	}
	if street == "" {
		return AddAddressCommand{}, errs.NewValueIsRequiredError("street")
	}
	if house == "" {
		return AddAddressCommand{}, errs.NewValueIsRequiredError("house")
	}
	if apartment == "" {
		return AddAddressCommand{}, errs.NewValueIsRequiredError("apartment")
	}

	return AddAddressCommand{
		basketID:  basketID,
		country:   country,
		city:      city,
		street:    street,
		house:     house,
		apartment: apartment,

		isValid: true,
	}, nil
}

func (c AddAddressCommand) IsValid() bool {
	return c.isValid
}
