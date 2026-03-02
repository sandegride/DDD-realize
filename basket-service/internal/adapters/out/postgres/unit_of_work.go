package postgres

import (
	"basket-service/internal/adapters/out/postgres/basketrepo"
	"basket-service/internal/adapters/out/postgres/goodrepo"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"basket-service/internal/pkg/outbox"
	"context"
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
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
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	committed := false
	defer func() {
		if !committed {
			if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
				log.Error(err)
			}
			u.clearTx()
		}
	}()

	if err := u.persistDomainEvents(ctx, u.tx); err != nil {
		return err
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	committed = true
	u.clearTx()

	return nil
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
}

func (u *UnitOfWork) persistDomainEvents(ctx context.Context, tx *gorm.DB) error {
	for _, agg := range u.trackedAggregates {
		outboxEvents, err := outbox.EncodeDomainEvents(agg.GetDomainEvents())
		if err != nil {
			return err
		}
		if len(outboxEvents) > 0 {
			if err := tx.WithContext(ctx).Create(&outboxEvents).Error; err != nil {
				return err
			}
		}
		agg.ClearDomainEvents()
	}
	return nil
}
