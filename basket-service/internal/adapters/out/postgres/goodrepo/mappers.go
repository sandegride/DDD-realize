package goodrepo

import (
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/core/domain/model/kernel"
)

func DomainToDTO(aggregate *good.Good) GoodDTO {
	var goodDTO GoodDTO
	goodDTO.ID = aggregate.ID()
	goodDTO.Title = aggregate.Title()
	goodDTO.Description = aggregate.Description()
	goodDTO.Price = aggregate.Price()
	goodDTO.Quantity = aggregate.Quantity()
	goodDTO.Weight = WeightDTO{
		Value: aggregate.Weight().Value(),
	}
	return goodDTO
}

func DtoToDomain(dto GoodDTO) *good.Good {
	weight, _ := kernel.NewWeight(dto.Weight.Value)
	aggregate := good.RestoreGood(dto.ID, dto.Title, dto.Description, dto.Price, dto.Quantity, weight)
	return aggregate
}
