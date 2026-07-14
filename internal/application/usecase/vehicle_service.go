package usecase

import (
	"context"

	vehicledto "github.com/Vanak/Backend/internal/application/dto/vehicle"
)

type VehicleService interface {
	RegisterVehicle(ctx context.Context, vehicleInfo vehicledto.VehicleDTORequest) error
	GetVehicle(vehicleHashedId string, userId uint) (*vehicledto.VehicleDTOResponse, error)
	DeleteVehicle(vehicleHashedId string, userId uint) error
	GetAllVehicles(userId uint) ([]vehicledto.VehicleListDTOResponse, error)
}
