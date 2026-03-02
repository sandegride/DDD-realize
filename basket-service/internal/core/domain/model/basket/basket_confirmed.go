package basket

import (
	"basket-service/internal/pkg/ddd"
	"github.com/google/uuid"
	"reflect"
)

const ConfirmedDomainEventName = "basket.confirmed.event"

var _ ddd.DomainEvent = &ConfirmedDomainEvent{}

type ConfirmedDomainEvent struct {
	ID   uuid.UUID
	Name string

	Payload BasketDTO

	isValid bool
}

func (e ConfirmedDomainEvent) GetID() uuid.UUID {
	return e.ID
}

func (e ConfirmedDomainEvent) GetName() string {
	return e.Name
}

func NewConfirmedDomainEvent(aggregate *Basket) *ConfirmedDomainEvent {
	dto := DomainToDTO(aggregate)
	domainEvent := ConfirmedDomainEvent{
		ID: uuid.New(),

		Payload: dto,

		isValid: true,
	}
	domainEvent.Name = reflect.TypeOf(domainEvent).Name()
	return &domainEvent
}

func NewEmptyConfirmedDomainEvent() ddd.DomainEvent {
	domainEvent := ConfirmedDomainEvent{}
	domainEvent.Name = reflect.TypeOf(domainEvent).Name()
	return &domainEvent
}

func (e *ConfirmedDomainEvent) IsValid() bool {
	return !e.isValid
}

type BasketDTO struct {
	ID             uuid.UUID
	BuyerID        uuid.UUID
	Address        AddressDTO
	DeliveryPeriod DeliveryPeriodDTO
	Items          []*ItemDTO
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

func DomainToDTO(aggregate *Basket) BasketDTO {
	var basketDTO BasketDTO
	basketDTO.ID = aggregate.ID()
	basketDTO.BuyerID = aggregate.BuyerID()
	basketDTO.Address = AddressDTO{
		Country:   aggregate.Address().Country(),
		City:      aggregate.Address().City(),
		Street:    aggregate.Address().Street(),
		House:     aggregate.Address().House(),
		Apartment: aggregate.Address().Apartment(),
	}

	if aggregate.DeliveryPeriod() != nil {
		basketDTO.DeliveryPeriod = DeliveryPeriodDTO{
			ID:   aggregate.DeliveryPeriod().ID(),
			Name: aggregate.DeliveryPeriod().Name(),
			From: aggregate.DeliveryPeriod().From(),
			To:   aggregate.DeliveryPeriod().To(),
		}
	}

	basketDTO.Items = make([]*ItemDTO, 0)
	for _, item := range aggregate.Items() {
		itemDTO := &ItemDTO{
			ID:          item.ID(),
			GoodID:      item.GoodID(),
			Title:       item.Title(),
			Description: item.Description(),
			Price:       item.Price(),
			Quantity:    item.Quantity(),
			BasketID:    aggregate.ID(),
		}
		basketDTO.Items = append(basketDTO.Items, itemDTO)
	}
	return basketDTO
}
