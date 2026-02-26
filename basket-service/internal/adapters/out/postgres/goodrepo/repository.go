package goodrepo

import (
	"basket-service/internal/adapters/out/postgres/shared"
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.GoodRepository = &Repository{}

type Repository struct {
	tracker shared.Tracker
}

func NewRepository(tracker shared.Tracker) (ports.GoodRepository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &Repository{
		tracker: tracker,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *good.Good) error {
	r.tracker.Track(aggregate)

	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *good.Good) error {
	r.tracker.Track(aggregate)

	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*good.Good, error) {
	dto := GoodDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, ID)
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError(ID.String(), nil)
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}
