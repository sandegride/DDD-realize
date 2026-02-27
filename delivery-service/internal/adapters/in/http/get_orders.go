package http

import (
	"delivery-service/internal/adapters/in/http/problems"
	"delivery-service/internal/core/application/usecases/queries"
	"delivery-service/internal/generated/servers"
	"delivery-service/internal/pkg/errs"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetOrders(c echo.Context) error {
	query, err := queries.NewGetNotCompletedOrdersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getNotCompletedOrdersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	var httpResponse = make([]servers.Order, 0, len(queryResponse.Orders))
	for _, courier := range queryResponse.Orders {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var courier = servers.Order{
			Id:       courier.ID,
			Location: location,
		}
		httpResponse = append(httpResponse, courier)
	}
	return c.JSON(http.StatusOK, httpResponse)
}
