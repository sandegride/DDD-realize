package courier

import (
	"delivery-service/internal/pkg/ddd"
	"delivery-service/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

type StoragePlace struct {
	baseEntity  *ddd.BaseEntity[uuid.UUID]
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

var (
	ErrTotalVolumeIsZeroOrLess            = errors.New("totalVolume should not be 0 or less")
	ErrStoragePlaceIsOccupied             = errors.New("storage place is occupied")
	ErrVolumeIsHigherThenTotalVolume      = errors.New("volume is higher than total volume")
	ErrCannotStoreOrderInThisStoragePlace = errors.New("cannot store order in this storage place")
	ErrOrderNotStoredInThisPlace          = errors.New("order is not stored in this place")
)

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if totalVolume <= 0 {
		return nil, ErrTotalVolumeIsZeroOrLess
	}

	return &StoragePlace{
		baseEntity:  ddd.NewBaseEntity(uuid.New()),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func RestoreStoragePlace(id uuid.UUID, name string, totalVolume int, orderID *uuid.UUID) *StoragePlace {
	return &StoragePlace{
		baseEntity:  ddd.NewBaseEntity(id),
		name:        name,
		totalVolume: totalVolume,
		orderID:     orderID,
	}
}

func (s *StoragePlace) Equal(other *StoragePlace) bool {
	if other == nil {
		return false
	}
	return s.baseEntity.Equal(other.baseEntity)
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.baseEntity.ID()
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume <= 0 {
		return false, errs.NewValueIsInvalidError("volume")
	}
	if s.isOccupied() {
		return false, nil
	}
	if volume > s.totalVolume {
		return false, nil
	}
	return true, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsRequiredError("orderID")
	}
	if volume <= 0 {
		return errs.NewValueIsRequiredError("volume")
	}

	canStore, err := s.CanStore(volume)
	if err != nil {
		return err
	}

	if !canStore {
		return ErrCannotStoreOrderInThisStoragePlace
	}

	s.orderID = &orderID
	return nil
}

func (sp *StoragePlace) Clear(orderID uuid.UUID) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsRequiredError("orderID")
	}
	if sp.orderID == nil || *sp.orderID != orderID {
		return ErrOrderNotStoredInThisPlace
	}

	sp.orderID = nil
	return nil
}

func (sp *StoragePlace) isOccupied() bool {
	return sp.orderID != nil
}
