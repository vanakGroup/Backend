package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type AddressRepository interface {
	CreateProvince(province *entities.Province) error
	CreateCity(city *entities.City) error
	GetProvinceByName(name string) (*entities.Province, error)
	GetCityByName(name string) (*entities.City, error)
}
