package http

import (
	"basket-service/internal/adapters/in/http/problems"
	"basket-service/internal/core/application/usecases/commands"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s *Server) Checkout(c echo.Context, basketID openapi_types.UUID) error {
	command, err := commands.NewCheckoutCommand(basketID)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.checkoutCommandHandler.Handle(c.Request().Context(), command)
	if err != nil {
		return problems.NewConflict(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}
