package ports

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/kernel"
	"context"
)

type DiscountClient interface {
	GetDiscount(ctx context.Context, basket *basket.Basket) (kernel.Discount, error)
	Close() error
}
