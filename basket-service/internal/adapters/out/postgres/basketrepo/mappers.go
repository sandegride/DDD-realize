package basketrepo

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/kernel"
)

func DomainToDTO(aggregate *basket.Basket) BasketDTO {
	var basketDTO BasketDTO
	basketDTO.ID = aggregate.ID()
	basketDTO.BuyerID = aggregate.BuyerID()
	basketDTO.Address = AddressDTO{
		Country:   aggregate.Address().Country(),
		City:      aggregate.Address().City(),
		Street:    aggregate.Address().Street(),
		House:     aggregate.Address().House(),
		Apartment: aggregate.Address().Apartment(),
	}

	if aggregate.DeliveryPeriod() != nil {
		id := aggregate.DeliveryPeriod().ID()
		basketDTO.DeliveryPeriodId = &id
	}

	basketDTO.Items = make([]*ItemDTO, 0)
	for _, item := range aggregate.Items() {
		itemDTO := &ItemDTO{
			ID:          item.ID(),
			GoodID:      item.GoodID(),
			Title:       item.Title(),
			Description: item.Description(),
			Price:       item.Price(),
			Quantity:    item.Quantity(),
			BasketID:    aggregate.ID(),
		}
		basketDTO.Items = append(basketDTO.Items, itemDTO)
	}
	basketDTO.Status = aggregate.Status()
	basketDTO.Total = aggregate.Total()
	return basketDTO
}

func DtoToDomain(dto BasketDTO) *basket.Basket {
	var aggregate *basket.Basket
	address, _ := kernel.NewAddress(dto.Address.Country, dto.Address.City, dto.Address.Street, dto.Address.House, dto.Address.Apartment)
	var deliveryPeriod *basket.DeliveryPeriod
	if dto.DeliveryPeriodId != nil {
		deliveryPeriod, _ = basket.GetDeliveryPeriodByID(*dto.DeliveryPeriodId)
	}

	var items []*basket.Item
	for _, dtoItem := range dto.Items {
		item := basket.RestoreItem(dtoItem.ID, dtoItem.GoodID, dtoItem.Title, dtoItem.Description, dtoItem.Price, dtoItem.Quantity)
		items = append(items, item)
	}
	aggregate = basket.RestoreBasket(dto.ID, dto.BuyerID, address, deliveryPeriod, items, dto.Status, dto.Total)
	return aggregate
}
