package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type ParkingRepository interface {
	GetAttributeByName(name string) (*entities.ParkingAttribute, error)
	CreateAttribute(parkingAttribute *entities.ParkingAttribute) error
	CreateParking(parking *entities.Parking) error
	GetAttributesByIDs(ids []uint) ([]*entities.ParkingAttribute, error)
	AddAttributes(parking *entities.Parking, attributes []*entities.ParkingAttribute) error
	CreateSlots(slots []entities.AvailabilitySlot) error
	UpdateParking(parking *entities.Parking) error
	CreateParkingImages(images []entities.ParkingImage) error
	FindByIDAndUserID(id uint, userID uint) (*entities.Parking, error)
	Delete(parking *entities.Parking) error
	GetByUserID(userID uint) ([]entities.Parking, error)
}
