package basket

import (
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrQuantityIsZeroOrLess = errors.New("quantity should not be 0 or less")
)

type Item struct {
	baseEntity *ddd.BaseEntity[uuid.UUID]

	goodID      uuid.UUID
	title       string
	description string
	price       float64
	quantity    int
}

func NewItem(good *good.Good, quantity int) (*Item, error) {
	if good == nil {
		return nil, errs.NewValueIsRequiredError("good")
	}
	if quantity <= 0 {
		return nil, ErrQuantityIsZeroOrLess
	}

	return &Item{
		baseEntity:  ddd.NewBaseEntity(uuid.New()),
		goodID:      good.ID(),
		title:       good.Title(),
		description: good.Description(),
		price:       good.Price(),
		quantity:    quantity,
	}, nil
}

func RestoreItem(id uuid.UUID, goodID uuid.UUID, title string, description string,
	price float64, quantity int) *Item {
	return &Item{
		baseEntity:  ddd.NewBaseEntity(id),
		goodID:      goodID,
		title:       title,
		description: description,
		price:       price,
		quantity:    quantity,
	}
}

func (i *Item) ID() uuid.UUID {
	return i.baseEntity.ID()
}

func (i *Item) GoodID() uuid.UUID {
	return i.goodID
}

func (i *Item) Title() string {
	return i.title
}

func (i *Item) Description() string {
	return i.description
}

func (i *Item) Price() float64 {
	return i.price
}

func (i *Item) Quantity() int {
	return i.quantity
}

func (i *Item) Equal(other *Item) bool {
	if other == nil {
		return false
	}
	return i.baseEntity.Equal(other.baseEntity)
}

func (i *Item) setQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityIsZeroOrLess
	}
	i.quantity = quantity
	return nil
}

func (i *Item) getTotal() float64 {
	total := float64(i.quantity) * i.price
	return total
}
