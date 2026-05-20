package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type ParkingRepository interface {
	GetAttributeByName(name string) (*entities.ParkingAttribute, error)
	CreateAttribute(parkingAttribute *entities.ParkingAttribute) error
	CreateParking(parking *entities.Parking) error
	GetAttributesByIDs(ids []uint) ([]*entities.ParkingAttribute, error)
	AddAttributes(parking *entities.Parking, attributes []*entities.ParkingAttribute) error
}
