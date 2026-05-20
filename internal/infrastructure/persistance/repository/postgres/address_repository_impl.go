package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type AddressRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewAddressRepository(cfg *bootstrap.Config, db *gorm.DB) *AddressRepository {
	return &AddressRepository{
		cfg: cfg,
		db:  db,
	}
}

func (ar *AddressRepository) CreateProvince(province *entities.Province) error {
	return ar.db.Create(province).Error
}

func (ar *AddressRepository) CreateCity(city *entities.City) error {
	return ar.db.Create(city).Error
}

func (ar *AddressRepository) GetProvinceByName(name string) (*entities.Province, error) {
	var province entities.Province
	result := ar.db.Where("name = ?", name).First(&province)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &province, nil
}

func (ar *AddressRepository) GetCityByName(name string) (*entities.City, error) {
	var city entities.City
	result := ar.db.Where("name = ?", name).First(&city)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &city, nil
}
