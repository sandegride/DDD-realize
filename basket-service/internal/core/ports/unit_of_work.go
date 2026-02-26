package ports

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	BasketRepository() BasketRepository
	GoodRepository() GoodRepository
	RollbackUnlessCommitted(ctx context.Context)
}
