package postgres

import (
	"basket-service/internal/adapters/out/postgres/basketrepo"
	"basket-service/internal/adapters/out/postgres/goodrepo"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"context"
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	committed         bool
	trackedAggregates []ddd.AggregateRoot
	basketRepository  ports.BasketRepository
	goodRepository    ports.GoodRepository
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow := &UnitOfWork{
		db: db,
	}

	goodRepo, err := goodrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.goodRepository = goodRepo

	basketRepo, err := basketrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.basketRepository = basketRepo

	return uow, nil
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) GoodRepository() ports.GoodRepository {
	return u.goodRepository
}

func (u *UnitOfWork) BasketRepository() ports.BasketRepository {
	return u.basketRepository
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
	u.committed = false
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	u.committed = true
	u.clearTx()
	return nil
}

func (u *UnitOfWork) RollbackUnlessCommitted(ctx context.Context) {
	if u.tx != nil && !u.committed {
		if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
			log.Error(err)
		}
		u.clearTx()
	}
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
	u.committed = false
}
