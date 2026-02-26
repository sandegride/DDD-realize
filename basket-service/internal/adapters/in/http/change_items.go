package http

import (
	"basket-service/internal/adapters/in/http/problems"
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/generated/servers"
	"basket-service/internal/pkg/errs"
	"errors"
	"github.com/labstack/echo/v4"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s *Server) ChangeItems(c echo.Context, basketID openapitypes.UUID) error {
	var item servers.Item
	err := c.Bind(&item)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	command, err := commands.NewChangeItemsCommand(basketID, basketID, item.GoodId, int(item.Quantity))
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.changeItemsCommandHandler.Handle(c.Request().Context(), command)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
		return c.JSON(http.StatusConflict, problems.NewConflict(err.Error(), err.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}
