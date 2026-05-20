package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type ParkingRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewParkingRepository(cfg *bootstrap.Config, db *gorm.DB) *ParkingRepository {
	return &ParkingRepository{
		cfg: cfg,
		db:  db,
	}
}

func (ps *ParkingRepository) GetAttributeByName(name string) (*entities.ParkingAttribute, error) {
	var attribute entities.ParkingAttribute
	result := ps.db.Where("name = ?", name).First(&attribute)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &attribute, nil
}

func (ps *ParkingRepository) CreateAttribute(parkingAttribute *entities.ParkingAttribute) error {
	return ps.db.Create(parkingAttribute).Error
}

func (ps *ParkingRepository) CreateParking(parking *entities.Parking) error {
	return ps.db.Create(parking).Error
}

func (ps *ParkingRepository) GetAttributesByIDs(ids []uint) ([]*entities.ParkingAttribute, error) {
	var attributes []*entities.ParkingAttribute
	result := ps.db.Where("id IN ?", ids).Find(&attributes)
	if result.Error != nil {
		return nil, result.Error
	}
	return attributes, nil
}

func (ps *ParkingRepository) AddAttributes(parking *entities.Parking, attributes []*entities.ParkingAttribute) error {
	return ps.db.Model(parking).Association("Attributes").Append(attributes)
}
