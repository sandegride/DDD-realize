package ports

import (
	"context"
	"delivery-service/internal/pkg/ddd"
)

type OrderProducer interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
	Close() error
}
