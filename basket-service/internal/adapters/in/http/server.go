package http

import (
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/core/application/usecases/queries"
	"basket-service/internal/pkg/errs"
)

type Server struct {
	addAddressCommandHandler        commands.AddAddressCommandHandler
	addDeliveryPeriodCommandHandler commands.AddDeliveryPeriodCommandHandler
	changeItemsCommandHandler       commands.ChangeItemsCommandHandler
	checkoutCommandHandler          commands.CheckoutCommandHandler
	getBasketQueryHandler           queries.GetBasketQueryHandler
}

func NewServer(
	addAddressCommandHandler commands.AddAddressCommandHandler,
	addDeliveryPeriodCommandHandler commands.AddDeliveryPeriodCommandHandler,
	changeItemsCommandHandler commands.ChangeItemsCommandHandler,
	checkoutCommandHandler commands.CheckoutCommandHandler,
	getBasketQueryHandler queries.GetBasketQueryHandler,
) (*Server, error) {
	if addAddressCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("addAddressCommandHandler")
	}
	if addDeliveryPeriodCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("addDeliveryPeriodCommandHandler")
	}
	if changeItemsCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("changeItemsCommandHandler")
	}
	if checkoutCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("checkoutCommandHandler")
	}
	if getBasketQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getBasketQueryHandler")
	}
	return &Server{
		addAddressCommandHandler:        addAddressCommandHandler,
		addDeliveryPeriodCommandHandler: addDeliveryPeriodCommandHandler,
		changeItemsCommandHandler:       changeItemsCommandHandler,
		checkoutCommandHandler:          checkoutCommandHandler,
		getBasketQueryHandler:           getBasketQueryHandler,
	}, nil
}
