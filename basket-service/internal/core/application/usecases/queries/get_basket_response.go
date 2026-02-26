package queries

import (
	"github.com/google/uuid"
)

type GetBasketResponse struct {
	ID      uuid.UUID
	Address Address
	Status  Status
}

type Status string

type Address struct {
	Country   string
	City      string
	Street    string
	House     string
	Apartment string
}
