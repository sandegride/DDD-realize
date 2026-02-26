package ports

import (
	"basket-service/internal/core/domain/model/good"
	"context"
	"github.com/google/uuid"
)

type GoodRepository interface {
	Add(ctx context.Context, aggregate *good.Good) error
	Update(ctx context.Context, aggregate *good.Good) error
	Get(ctx context.Context, ID uuid.UUID) (*good.Good, error)
}
