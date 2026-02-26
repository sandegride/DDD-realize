package kernel

import "basket-service/internal/pkg/errs"

type Address struct {
	country   string
	city      string
	street    string
	house     string
	apartment string

	valid bool
}

func NewAddress(country, city, street, house, apartment string) (Address, error) {
	if country == "" {
		return Address{}, errs.NewValueIsRequiredError("country")
	}
	if city == "" {
		return Address{}, errs.NewValueIsRequiredError("city")
	}
	if street == "" {
		return Address{}, errs.NewValueIsRequiredError("street")
	}
	if house == "" {
		return Address{}, errs.NewValueIsRequiredError("house")
	}
	if apartment == "" {
		return Address{}, errs.NewValueIsRequiredError("apartment")
	}

	return Address{country, city, street, house, apartment, true}, nil
}

func (a Address) Equal(other Address) bool {
	return a.country == other.country && a.city == other.city && a.street == other.street && a.house == other.house && a.apartment == other.apartment
}

func (a Address) IsValid() bool {
	return a.valid
}

func (a Address) Country() string {
	return a.country
}

func (a Address) City() string {
	return a.city
}

func (a Address) Street() string {
	return a.street
}

func (a Address) House() string {
	return a.house
}

func (a Address) Apartment() string {
	return a.apartment
}
