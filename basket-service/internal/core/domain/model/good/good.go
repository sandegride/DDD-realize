package good

import (
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
)

var (
	Goods = []*Good{
		Bread(),
		Milk(),
		Eggs(),
		Sausage(),
		Coffee(),
		Sugar(),
		Gum(),
		Candy(),
		Snack(),
	}
)

type Good struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]

	title       string
	description string
	price       float64
	quantity    int
	weight      kernel.Weight
}

func NewGood(id uuid.UUID, title string, description string, price float64, quantity int, weight kernel.Weight) (*Good, error) {
	if id == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("id")
	}
	if title == "" {
		return nil, errs.NewValueIsRequiredError("title")
	}
	if description == "" {
		return nil, errs.NewValueIsRequiredError("description")
	}
	if price < 0 {
		return nil, errs.NewValueIsInvalidError("price")
	}
	if quantity <= 0 {
		return nil, errs.NewValueIsInvalidError("quantity")
	}
	if !weight.IsValid() {
		return nil, errs.NewValueIsRequiredError("weight")
	}

	return &Good{
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),

		title:       title,
		description: description,
		price:       price,
		quantity:    quantity,
		weight:      weight,
	}, nil
}

func RestoreGood(id uuid.UUID, title string, description string, price float64, quantity int, weight kernel.Weight) *Good {
	return &Good{
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),

		title:       title,
		description: description,
		price:       price,
		quantity:    quantity,
		weight:      weight,
	}
}

func (g *Good) ID() uuid.UUID {
	return g.baseAggregate.ID()
}

func (g *Good) Equal(other *Good) bool {
	if other == nil {
		return false
	}
	return g.baseAggregate.Equal(other.baseAggregate)
}

func (g *Good) ClearDomainEvents() {
	g.baseAggregate.ClearDomainEvents()
}

func (g *Good) GetDomainEvents() []ddd.DomainEvent {
	return g.baseAggregate.GetDomainEvents()
}

func (g *Good) RaiseDomainEvent(event ddd.DomainEvent) {
	g.baseAggregate.RaiseDomainEvent(event)
}

func (g *Good) Title() string {
	return g.title
}

func (g *Good) Description() string {
	return g.description
}

func (g *Good) Price() float64 {
	return g.price
}

func (g *Good) Quantity() int {
	return g.quantity
}

func (g *Good) Weight() kernel.Weight {
	return g.weight
}

func (g *Good) ChangeStocks(quantity int) error {
	if quantity < 0 {
		return errs.NewValueIsInvalidError("quantity")
	}
	g.quantity = quantity
	return nil
}

func Bread() *Good {
	weight, _ := kernel.NewWeight(6)
	good, _ := NewGood(
		uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c"),
		"Хлеб",
		"Описание хлеба",
		100,
		100,
		weight)
	return good
}

func Milk() *Good {
	weight, _ := kernel.NewWeight(9)
	good, _ := NewGood(
		uuid.MustParse("e8cb8a0b-d302-485a-801c-5fb50aced4d5"),
		"Молоко",
		"Описание молока",
		200,
		100,
		weight)
	return good
}

func Eggs() *Good {
	weight, _ := kernel.NewWeight(8)
	good, _ := NewGood(
		uuid.MustParse("a1d48be9-4c98-4371-97c0-064bde03874e"),
		"Яйца",
		"Описание яиц",
		300,
		100,
		weight)
	return good
}

func Sausage() *Good {
	weight, _ := kernel.NewWeight(4)
	good, _ := NewGood(
		uuid.MustParse("34b1e64a-6471-44a0-8c4a-e5d21584a76c"),
		"Колбаса",
		"Описание колбасы",
		400,
		100,
		weight)
	return good
}

func Coffee() *Good {
	weight, _ := kernel.NewWeight(7)
	good, _ := NewGood(
		uuid.MustParse("292dc3c5-2bdd-4e0c-bd75-c5e8b07a8792"),
		"Кофе",
		"Описание кофе",
		500,
		100,
		weight)
	return good
}

func Sugar() *Good {
	weight, _ := kernel.NewWeight(1)
	good, _ := NewGood(
		uuid.MustParse("a3fcc8e1-d2a3-4bd6-9421-c82019e21c2d"),
		"Сахар",
		"Описание сахара",
		600,
		100,
		weight)
	return good
}

func Gum() *Good {
	weight, _ := kernel.NewWeight(1)
	good, _ := NewGood(
		uuid.MustParse("3ecc7f05-7081-4155-9ca1-a3e371faa661"),
		"Жвачка",
		"Промо товар",
		1,
		100,
		weight)
	return good
}

func Candy() *Good {
	weight, _ := kernel.NewWeight(1)
	good, _ := NewGood(
		uuid.MustParse("67dda004-f324-4868-af9a-88f77d1d28fc"),
		"Конфета",
		"Промо товар",
		1,
		100,
		weight)
	return good
}

func Snack() *Good {
	weight, _ := kernel.NewWeight(1)
	good, _ := NewGood(
		uuid.MustParse("5399a4c3-8b32-418a-b477-cedb6be815d0"),
		"Перекус",
		"Промо товар",
		1,
		100,
		weight)
	return good
}
