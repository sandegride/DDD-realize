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

func (s *Server) GetCouriers(c echo.Context) error {
	query, err := queries.NewGetAllCouriersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	queryResponse, err := s.getAllCouriersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	var httpResponse = make([]servers.Courier, 0, len(queryResponse.Couriers))
	for _, courier := range queryResponse.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var courier = servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}
		httpResponse = append(httpResponse, courier)
	}
	return c.JSON(http.StatusOK, httpResponse)
}
