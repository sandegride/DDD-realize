package ports

import (
	"basket-service/internal/pkg/ddd"
	"context"
)

type BasketProducer interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
	Close() error
}
