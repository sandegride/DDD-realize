package basket

import (
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrHasAlreadyBeenIssuedError           = errors.New("basket has already been issued")
	ErrItemsShouldNotBeEmptyError          = errors.New("items should not empty")
	ErrAddressShouldNotBeEmptyError        = errors.New("address should not be empty")
	ErrDeliveryPeriodShouldNotBeEmptyError = errors.New("delivery period should not be empty")
)

type Basket struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]

	buyerID        uuid.UUID
	address        kernel.Address
	deliveryPeriod *DeliveryPeriod
	items          []*Item
	status         Status
	total          float64
}

func NewBasket(buyerID uuid.UUID) (*Basket, error) {
	if buyerID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("buyerID")
	}

	return &Basket{
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](uuid.New()),

		buyerID: buyerID,
		status:  StatusCreated,
		items:   make([]*Item, 0),
	}, nil
}

func RestoreBasket(id uuid.UUID, buyerID uuid.UUID, address kernel.Address, deliveryPeriod *DeliveryPeriod,
	items []*Item, status Status, total float64) *Basket {

	if items == nil {
		items = []*Item{}
	}
	return &Basket{
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),

		buyerID:        buyerID,
		address:        address,
		deliveryPeriod: deliveryPeriod,
		items:          items,
		status:         status,
		total:          total,
	}
}

func (b *Basket) ID() uuid.UUID {
	return b.baseAggregate.ID()
}

func (b *Basket) Equal(other *Basket) bool {
	if other == nil {
		return false
	}
	return b.baseAggregate.Equal(other.baseAggregate)
}

func (b *Basket) ClearDomainEvents() {
	b.baseAggregate.ClearDomainEvents()
}

func (b *Basket) GetDomainEvents() []ddd.DomainEvent {
	return b.baseAggregate.GetDomainEvents()
}

func (b *Basket) RaiseDomainEvent(event ddd.DomainEvent) {
	b.baseAggregate.RaiseDomainEvent(event)
}

func (b *Basket) BuyerID() uuid.UUID {
	return b.buyerID
}

func (b *Basket) Address() kernel.Address {
	return b.address
}

func (b *Basket) DeliveryPeriod() *DeliveryPeriod {
	return b.deliveryPeriod
}

func (b *Basket) Items() []*Item {
	return b.items
}

func (b *Basket) Status() Status {
	return b.status
}

func (b *Basket) Total() float64 {
	return b.total
}

func (b *Basket) Change(good *good.Good, quantity int) error {
	if b.status == StatusConfirmed {
		return ErrHasAlreadyBeenIssuedError
	}
	if good == nil {
		return errs.NewValueIsRequiredError("good")
	}
	if quantity < 0 {
		return errs.NewValueIsRequiredError("quantity")
	}

	item, _, err := b.FindItemByGoodID(good.ID())
	if err != nil {
		return err
	}

	if item != nil {
		if quantity == 0 {
			err := b.removeItem(item)
			if err != nil {
				return err
			}
		} else {
			err := item.setQuantity(quantity)
			if err != nil {
				return err
			}
		}
	} else {
		var item, err = NewItem(good, quantity)
		if err != nil {
			return err
		}
		err = b.addItem(item)
		if err != nil {
			return err
		}
	}

	b.total = b.getTotal()
	return nil
}

func (b *Basket) AddAddress(address kernel.Address) error {
	if !address.IsValid() {
		return errs.NewValueIsRequiredError("address")
	}
	if b.status == StatusConfirmed {
		return ErrHasAlreadyBeenIssuedError
	}

	b.address = address
	return nil
}

func (b *Basket) AddDeliveryPeriod(deliveryPeriod *DeliveryPeriod) error {
	if deliveryPeriod == nil {
		return errs.NewValueIsRequiredError("deliveryPeriod")
	}
	if b.status == StatusConfirmed {
		return ErrHasAlreadyBeenIssuedError
	}

	b.deliveryPeriod = deliveryPeriod
	return nil
}

func (b *Basket) Checkout(discount kernel.Discount) error {
	if !b.status.IsValid() {
		return errs.NewValueIsRequiredError("status")
	}
	if b.status == StatusConfirmed {
		return ErrHasAlreadyBeenIssuedError
	}
	if len(b.items) <= 0 {
		return ErrItemsShouldNotBeEmptyError
	}
	if !b.address.IsValid() {
		return ErrAddressShouldNotBeEmptyError
	}
	if b.deliveryPeriod == nil {
		return ErrDeliveryPeriodShouldNotBeEmptyError
	}
	if !discount.IsValid() {
		return errs.NewValueIsRequiredError("discount")
	}

	// Рассчитываем итоговую стоимость, учитывая размер скидки
	var total = b.getTotal()
	totalWithDiscount, err := discount.Apply(total)
	if err != nil {
		return err
	}
	b.total = totalWithDiscount

	// Меняем статус
	b.status = StatusConfirmed
	return nil
}

func (b *Basket) FindItemByGoodID(goodID uuid.UUID) (*Item, int, error) {
	if goodID == uuid.Nil {
		return nil, 0, errs.NewValueIsInvalidError("goodID")
	}

	for i := range b.items {
		if b.items[i].GoodID() == goodID {
			return b.items[i], i, nil
		}
	}
	return nil, 0, nil
}

func (b *Basket) addItem(item *Item) error {
	if item == nil {
		return errs.NewValueIsInvalidError("item")
	}
	b.items = append(b.items, item)
	return nil
}

func (b *Basket) removeItem(item *Item) error {
	if item == nil {
		return errs.NewValueIsInvalidError("item")
	}
	_, index, err := b.FindItemByGoodID(item.goodID)
	if err != nil {
		return err
	}
	b.items = append(b.items[:index], b.items[index+1:]...)
	return nil
}

func (b *Basket) getTotal() float64 {
	var total float64
	for i := range b.items {
		total += b.items[i].getTotal()
	}
	return total
}
