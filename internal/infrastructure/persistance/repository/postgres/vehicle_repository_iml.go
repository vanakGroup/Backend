package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type VehicleRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewVehicleRepository(cfg *bootstrap.Config, db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{
		cfg: cfg,
		db:  db,
	}
}

func (vr *VehicleRepository) FindVehiclesByLicensePlate(licensePlate string) bool {
	var exists bool

	vr.db.Model(&entities.Vehicle{}).
		Select("count(*) > 0").
		Where("license_plate = ?", licensePlate).
		Find(&exists)

	return exists
}

func (vr *VehicleRepository) Create(vehicle *entities.Vehicle) error {
	return vr.db.Create(vehicle).Error
}

func (vr *VehicleRepository) FindUserVehicle(vehicleId uint, userId uint) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle

	err := vr.db.Model(&entities.Vehicle{}).
		Where("user_id = ?", userId).
		Where("id = ? ", vehicleId).
		First(&vehicle).Error

	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (vr *VehicleRepository) RemoveVehicle(vehicleId uint, userId uint) error {
	err := vr.db.Where("id = ? AND user_id = ?", vehicleId, userId).
		Delete(&entities.Vehicle{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (vr *VehicleRepository) FindAllUserVehicles(userId uint) ([]entities.Vehicle, error) {
	var vehicles []entities.Vehicle

	err := vr.db.Where("user_id = ?", userId).
		Find(&vehicles).Error

	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (vr *VehicleRepository) FindVehicleByID(vehicleID uint) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle
	err := vr.db.Model(&entities.Vehicle{}).
		Where("id = ?", vehicleID).
		First(&vehicle).Error

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}
