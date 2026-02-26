package ports

import (
	"basket-service/internal/core/domain/model/basket"
	"context"
	"github.com/google/uuid"
)

type BasketRepository interface {
	Add(ctx context.Context, aggregate *basket.Basket) error
	Update(ctx context.Context, aggregate *basket.Basket) error
	Get(ctx context.Context, ID uuid.UUID) (*basket.Basket, error)
}
