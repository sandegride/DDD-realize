package courierrepo

import (
	"delivery-service/internal/core/domain/model/courier"
	"delivery-service/internal/core/domain/model/kernel"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO
	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed()
	courierDTO.Location = LocationDTO{
		X: aggregate.Location().X(),
		Y: aggregate.Location().Y(),
	}
	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	var aggregate *courier.Courier
	var storagePlaces []*courier.StoragePlace
	for _, dtoStoragePlace := range dto.StoragePlaces {
		item := courier.RestoreStoragePlace(dtoStoragePlace.ID, dtoStoragePlace.Name,
			dtoStoragePlace.TotalVolume, dtoStoragePlace.OrderID)
		storagePlaces = append(storagePlaces, item)
	}
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	aggregate = courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, location, storagePlaces)
	return aggregate
}
