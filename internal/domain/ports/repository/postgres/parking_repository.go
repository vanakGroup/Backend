package postgres

import (
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/filter"
)

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
	HasActiveReservations(parkingID uint) (bool, error)
	GetSlotsByParkingID(parkingID uint) ([]entities.AvailabilitySlot, error)
	DeleteSlotsWithIDs(slotIDs []uint) error
	DeleteParkingImagesWithIDs(imageIDs []uint) error
	UpdateAttributes(parking *entities.Parking, attributes []entities.ParkingAttribute) error
	GetDetailByIDAndUserID(parkingID uint, userID uint) (*entities.Parking, error)
	GetAllAttributes() ([]entities.ParkingAttribute, error)
	GetParkingByIDAndUserID(userID uint, parkingID uint) (*entities.Parking, error)
	GetAllByUserID(userID uint) ([]*entities.Parking, error)
	Search(f filter.ParkingFilter, options *filter.QueryOptions) ([]entities.Parking, error)
	Count(f filter.ParkingFilter) (int64, error)
	GetAll(options *filter.QueryOptions) ([]entities.Parking, error)
	CountAll() (int64, error)
	GetByID(id uint) (*entities.Parking, error)
}
