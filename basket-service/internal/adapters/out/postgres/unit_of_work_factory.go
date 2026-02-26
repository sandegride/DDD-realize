package postgres

import (
	"basket-service/internal/core/ports"
	"basket-service/internal/pkg/errs"
	"context"
	"gorm.io/gorm"
)

type unitOfWorkFactory struct {
	db *gorm.DB
}

func NewUnitOfWorkFactory(db *gorm.DB) (ports.UnitOfWorkFactory, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &unitOfWorkFactory{db: db}, nil
}

func (f *unitOfWorkFactory) New(ctx context.Context) (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db.WithContext(ctx))
}
