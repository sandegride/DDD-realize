package http

import (
	"basket-service/internal/adapters/in/http/problems"
	"basket-service/internal/core/application/usecases/queries"
	"basket-service/internal/generated/servers"
	"basket-service/internal/pkg/errs"
	"errors"
	"github.com/labstack/echo/v4"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func (s *Server) GetBasket(c echo.Context, basketID openapitypes.UUID) error {
	query, err := queries.NewGetBasketQuery(basketID)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getBasketQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	address := servers.Address{
		Country:   queryResponse.Address.Country,
		City:      queryResponse.Address.City,
		Street:    queryResponse.Address.Street,
		House:     queryResponse.Address.House,
		Apartment: queryResponse.Address.Apartment,
	}

	var httpResponse = servers.Basket{
		Id:             &queryResponse.ID,
		Items:          nil,
		Address:        &address,
		DeliveryPeriod: nil,
		Status:         nil,
	}

	return c.JSON(http.StatusOK, httpResponse)
}
