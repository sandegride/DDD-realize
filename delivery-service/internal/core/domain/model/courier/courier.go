package courier

import (
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/domain/model/order"
	"delivery-service/internal/pkg/ddd"
	"delivery-service/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"math"
)

type Courier struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

var (
	ErrNoSuitableStoragePlace = errors.New("no suitable storage place")
	ErrNotFoundOrder          = errors.New("order not found")
)

const BackpackName = "Сумка"
const BackpackSize = 10

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if speed <= 0 {
		return nil, errs.NewValueIsRequiredError("speed")
	}
	if !location.IsValid() {
		return nil, errs.NewValueIsRequiredError("location")
	}

	courier := &Courier{
		baseAggregate: ddd.NewBaseAggregate[uuid.UUID](uuid.New()),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: make([]*StoragePlace, 0),
	}

	err := courier.AddStoragePlace(BackpackName, BackpackSize)
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func RestoreCourier(id uuid.UUID, name string, speed int, location kernel.Location, places []*StoragePlace) *Courier {
	return &Courier{
		baseAggregate: ddd.NewBaseAggregate(id),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: places,
	}
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}
	return c.baseAggregate.ID() == other.baseAggregate.ID()
}

func (c *Courier) ClearDomainEvents() {
	c.baseAggregate.ClearDomainEvents()
}

func (c *Courier) GetDomainEvents() []ddd.DomainEvent {
	return c.baseAggregate.GetDomainEvents()
}

func (c *Courier) RaiseDomainEvent(event ddd.DomainEvent) {
	c.baseAggregate.RaiseDomainEvent(event)
}

func (c *Courier) ID() uuid.UUID {
	return c.baseAggregate.ID()
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []StoragePlace {
	res := make([]StoragePlace, len(c.storagePlaces))
	for i, storagePlace := range c.storagePlaces {
		res[i] = *storagePlace
	}
	return res
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	storagePlace, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}
	c.storagePlaces = append(c.storagePlaces, storagePlace)
	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, errs.NewValueIsRequiredError("order")
	}

	for _, storagePlace := range c.storagePlaces {
		canStore, err := storagePlace.CanStore(order.Volume())
		if err != nil {
			return false, err
		}

		if canStore {
			return true, nil
		}
	}
	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	canTakeOrder, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if !canTakeOrder {
		return ErrNoSuitableStoragePlace
	}

	for _, storagePlace := range c.storagePlaces {
		canStore, err := storagePlace.CanStore(order.Volume())
		if err != nil {
			return err
		}

		if canStore {
			err := storagePlace.Store(order.ID(), order.Volume())
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNoSuitableStoragePlace
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	storagePlace, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return errs.NewObjectNotFoundError("order", order.ID())
	}

	err = storagePlace.Clear(order.ID())
	if err != nil {
		return err
	}
	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsRequiredError("target")
	}
	distance, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}

	time := float64(distance) / float64(c.speed)
	return time, err
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsValid() {
		return errs.NewValueIsRequiredError("target")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	for _, storagePlace := range c.storagePlaces {
		if *storagePlace.OrderID() != uuid.Nil && *storagePlace.OrderID() == orderID {
			return storagePlace, nil
		}
	}
	return nil, nil
}

func (c *Courier) getBackpackOccupiedSize() int {
	occupiedSize := 0
	for _, storagePlace := range c.storagePlaces {
		occupiedSize += storagePlace.TotalVolume()
	}
	return occupiedSize
}
