package order

import (
	"delivery-service/internal/pkg/ddd"
	"github.com/google/uuid"
	"reflect"
)

var _ ddd.DomainEvent = &CompletedDomainEvent{}

type CompletedDomainEvent struct {
	// base
	ID   uuid.UUID
	Name string

	// payload
	OrderID     uuid.UUID
	OrderStatus string

	valid bool
}

func (e CompletedDomainEvent) GetID() uuid.UUID { return e.ID }

func (e CompletedDomainEvent) GetName() string {
	return e.Name
}

func NewCompletedDomainEvent(aggregate *Order) ddd.DomainEvent {
	domainEvent := CompletedDomainEvent{
		ID: uuid.New(),

		OrderID:     aggregate.ID(),
		OrderStatus: aggregate.Status().String(),

		valid: true,
	}
	domainEvent.Name = reflect.TypeOf(domainEvent).Name()
	return &domainEvent
}

func NewEmptyCompletedDomainEvent() ddd.DomainEvent {
	domainEvent := CompletedDomainEvent{}
	domainEvent.Name = reflect.TypeOf(domainEvent).Name()
	return &domainEvent
}

func (e CompletedDomainEvent) IsValid() bool {
	return !e.valid
}
