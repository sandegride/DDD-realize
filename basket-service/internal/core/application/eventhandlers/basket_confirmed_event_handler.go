package eventhandlers

import (
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"context"
)

type basketCompletedDomainEventHandler struct {
	basketProducer ports.BasketProducer
}

func NewBasketCompletedDomainEventHandler(
	basketProducer ports.BasketProducer) (ddd.EventHandler, error) {
	if basketProducer == nil {
		return nil, errs.NewValueIsRequiredError("basketProducer")
	}

	return &basketCompletedDomainEventHandler{basketProducer: basketProducer}, nil
}

func (eh *basketCompletedDomainEventHandler) Handle(ctx context.Context, domainEvent ddd.DomainEvent) error {
	err := eh.basketProducer.Publish(ctx, domainEvent)
	if err != nil {
		return err
	}
	return nil
}
