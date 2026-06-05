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

func (ar *AddressRepository) GetCityByID(id uint) (*entities.City, error) {
	var city entities.City

	result := ar.db.First(&city, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, result.Error
	}

	return &city, nil
}

func (ar *AddressRepository) CreateAddress(address *entities.Address) error {
	return ar.db.Create(address).Error
}

func (r *AddressRepository) GetAddressByParkingID(parkingID uint) (*entities.Address, error) {
	var address entities.Address
	err := r.db.Where("parking_id = ?", parkingID).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *AddressRepository) UpdateAddress(address *entities.Address) error {
	return r.db.Save(address).Error
}
