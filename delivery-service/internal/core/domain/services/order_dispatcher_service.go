package services

import (
	"delivery-service/internal/core/domain/model/courier"
	"delivery-service/internal/core/domain/model/order"
	"delivery-service/internal/pkg/errs"
	"errors"
	"math"
)

type OrderDispatcherService interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcherService = &orderDispatcherService{}

type orderDispatcherService struct{}

func NewOrderDispatcherService() OrderDispatcherService {
	return &orderDispatcherService{}
}

var (
	ErrOrderStatusIsNotCreatedStatus = errors.New("order status is not created status")
	ErrNotFoundFreeCourier           = errors.New("no free courier")
)

func (od *orderDispatcherService) Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}
	if !order.Status().IsCreated() {
		return nil, ErrOrderStatusIsNotCreatedStatus
	}
	if len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	bestCourier, err := od.findBestCourier(order, couriers)
	if err != nil {
		return nil, err
	}

	if err := bestCourier.TakeOrder(order); err != nil {
		return nil, err
	}
	if err := order.Assign(bestCourier.ID()); err != nil {
		return nil, err
	}

	return bestCourier, nil
}

func (od *orderDispatcherService) findBestCourier(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var bestCourier *courier.Courier
	minTime := math.MaxFloat64

	for _, c := range couriers {
		canTake, err := c.CanTakeOrder(order)
		if err != nil {
			continue
		}
		if !canTake {
			continue
		}

		time, err := c.CalculateTimeToLocation(order.Location())
		if err != nil {
			continue
		}

		if time < minTime {
			minTime = time
			bestCourier = c
		}
	}

	if bestCourier == nil {
		return nil, ErrNotFoundFreeCourier
	}

	return bestCourier, nil
}
