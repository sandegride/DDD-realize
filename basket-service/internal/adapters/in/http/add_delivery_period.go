package http

import (
	"basket-service/internal/adapters/in/http/problems"
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/generated/servers"
	"github.com/labstack/echo/v4"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s *Server) AddDeliveryPeriod(c echo.Context, basketID openapitypes.UUID) error {
	var deliveryPeriod servers.DeliveryPeriod
	err := c.Bind(&deliveryPeriod)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	command, err := commands.NewAddDeliveryPeriodCommand(basketID, commands.DeliveryPeriod(deliveryPeriod))
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.addDeliveryPeriodCommandHandler.Handle(c.Request().Context(), command)
	if err != nil {
		return problems.NewConflict(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}
