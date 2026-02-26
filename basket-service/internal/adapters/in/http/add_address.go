package http

import (
	"basket-service/internal/adapters/in/http/problems"
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/generated/servers"
	"github.com/labstack/echo/v4"
	openapitypes "github.com/oapi-codegen/runtime/types"

	"net/http"
)

func (s *Server) AddAddress(c echo.Context, basketID openapitypes.UUID) error {
	var address servers.Address
	err := c.Bind(&address)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	addAddressCommand, err := commands.NewAddAddressCommand(basketID, address.Country, address.City, address.Street, address.House, address.Apartment)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.addAddressCommandHandler.Handle(c.Request().Context(), addAddressCommand)
	if err != nil {
		return problems.NewConflict(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}
