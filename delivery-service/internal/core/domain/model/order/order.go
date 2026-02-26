package order

import (
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/pkg/ddd"
	"delivery-service/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrCannotCompleteNotAssignedOrder   = errors.New("can not complete not assigned order")
	ErrCannotAssignAlreadyAssignedOrder = errors.New("can not assign already assigned order")
)

type Order struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	courierID     *uuid.UUID

	location kernel.Location
	volume   int
	status   Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if !location.IsValid() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](orderID),
		location:      location,
		volume:        volume,
		status:        StatusCreated,
	}, nil
}

func RestoreOrder(id uuid.UUID, courierID *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		baseAggregate: ddd.NewBaseAggregate(id),
		courierID:     courierID,
		location:      location,
		volume:        volume,
		status:        status,
	}
}

func (o *Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}
	return o.baseAggregate.Equal(other.baseAggregate)
}

func (o *Order) ClearDomainEvents() {
	o.baseAggregate.ClearDomainEvents()
}

func (o *Order) GetDomainEvents() []ddd.DomainEvent {
	return o.baseAggregate.GetDomainEvents()
}

func (o *Order) RaiseDomainEvent(event ddd.DomainEvent) {
	o.baseAggregate.RaiseDomainEvent(event)
}

func (o *Order) ID() uuid.UUID {
	return o.baseAggregate.ID()
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if courierID == uuid.Nil {
		return errs.NewValueIsRequiredError("courierID")
	}
	if o.status != StatusCreated {
		return ErrCannotAssignAlreadyAssignedOrder
	}

	o.courierID = &courierID
	o.status = StatusAssigned
	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return ErrCannotCompleteNotAssignedOrder
	}
	o.status = StatusCompleted

	return nil
}
