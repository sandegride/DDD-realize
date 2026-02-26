package goodrepo

import "github.com/google/uuid"

type GoodDTO struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Title       string
	Description string
	Price       float64
	Quantity    int
	Weight      WeightDTO `gorm:"embedded;embeddedPrefix:weight_"`
}

type WeightDTO struct {
	Value int
}

func (GoodDTO) TableName() string {
	return "goods"
}
