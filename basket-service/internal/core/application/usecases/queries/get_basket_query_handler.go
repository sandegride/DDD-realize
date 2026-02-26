package queries

import (
	"basket-service/internal/pkg/errs"
	"gorm.io/gorm"
)

type GetBasketQueryHandler interface {
	Handle(GetBasketQuery) (GetBasketResponse, error)
}

type getBasketQueryHandler struct {
	db *gorm.DB
}

func NewGetBasketQueryHandler(db *gorm.DB) (GetBasketQueryHandler, error) {
	if db == nil {
		return &getBasketQueryHandler{}, errs.NewValueIsRequiredError("db")
	}
	return &getBasketQueryHandler{db: db}, nil
}

func (q *getBasketQueryHandler) Handle(query GetBasketQuery) (GetBasketResponse, error) {
	if query.IsValid() {
		return GetBasketResponse{}, errs.NewValueIsRequiredError("query")
	}

	var response GetBasketResponse
	result := q.db.Raw("SELECT * FROM baskets WHERE id = ?", query.BasketID()).Scan(&response)
	if result.RowsAffected == 0 {
		return GetBasketResponse{}, errs.NewObjectNotFoundError(query.BasketID().String(), nil)
	}
	if result.Error != nil {
		return GetBasketResponse{}, result.Error
	}
	return response, nil
}
