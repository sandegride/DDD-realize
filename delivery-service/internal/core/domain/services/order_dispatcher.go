package services

import (
	"delivery-service/internal/core/domain/model/courier"
	ord "delivery-service/internal/core/domain/model/order"
	"delivery-service/internal/pkg/errs"
	"errors"
	"math"
)

var (
	ErrSuitableCourierWasNotFound = errors.New("suitable courier was not found")
	ErrOrderAlreadyAssigned       = errors.New("order is already assigned")
)

type OrderDispatcher interface {
	Dispatch(order *ord.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcher = &orderDispatcher{}

type orderDispatcher struct {
}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (p *orderDispatcher) Dispatch(order *ord.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}
	if len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}
	if order.Status() != ord.StatusCreated {
		return nil, ErrOrderAlreadyAssigned
	}

	bestCourier, err := p.findBestCourier(order, couriers)
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

func (p *orderDispatcher) findBestCourier(order *ord.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var bestCourier *courier.Courier
	minTime := math.MaxFloat64

	for _, c := range couriers {
		canTake, err := c.CanTakeOrder(order)
		if err != nil {
			return nil, err
		}
		if !canTake {
			continue
		}

		time, err := c.CalculateTimeToLocation(order.Location())
		if err != nil {
			return nil, err
		}

		if time < minTime {
			minTime = time
			bestCourier = c
		}
	}

	if bestCourier == nil {
		return nil, ErrSuitableCourierWasNotFound
	}
	return bestCourier, nil
}
