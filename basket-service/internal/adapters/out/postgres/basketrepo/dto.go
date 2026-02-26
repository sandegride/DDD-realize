package basketrepo

import (
	"basket-service/internal/core/domain/model/basket"
	"github.com/google/uuid"
)

type BasketDTO struct {
	ID               uuid.UUID     `gorm:"type:uuid;primaryKey"`
	BuyerID          uuid.UUID     `gorm:"type:uuid;index"`
	Address          AddressDTO    `gorm:"embedded;embeddedPrefix:address_"`
	DeliveryPeriodId *int          `gorm:"index;default:null"`
	Items            []*ItemDTO    `gorm:"foreignKey:BasketID;constraint:OnDelete:CASCADE;"`
	Status           basket.Status `gorm:"type:varchar(20)"`
	Total            float64
}

type AddressDTO struct {
	Country   string
	City      string
	Street    string
	House     string
	Apartment string
}

type DeliveryPeriodDTO struct {
	ID   int
	Name string
	From int
	To   int
}

type ItemDTO struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	GoodID      uuid.UUID
	Title       string
	Description string
	Price       float64
	Quantity    int
	BasketID    uuid.UUID `gorm:"type:uuid;index"`
}

func (BasketDTO) TableName() string {
	return "baskets"
}

func (ItemDTO) TableName() string {
	return "items"
}

func (DeliveryPeriodDTO) TableName() string {
	return "delivery_periods"
}
