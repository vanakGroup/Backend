package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type VehicleRepository interface {
	FindVehiclesByLicensePlate(licensePlate string) bool
	Create(*entities.Vehicle) error
	FindUserVehicle(vehicleId uint, userId uint) (*entities.Vehicle, error)
	RemoveVehicle(vehicleId uint, userId uint) error
	FindAllUserVehicles(userId uint) ([]entities.Vehicle, error)
}
