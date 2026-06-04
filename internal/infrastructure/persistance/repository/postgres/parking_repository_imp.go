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

func (ps *ParkingRepository) CreateSlots(slots []entities.AvailabilitySlot) error {
	return ps.db.Create(&slots).Error
}

func (r *ParkingRepository) UpdateParking(parking *entities.Parking) error {
	return r.db.Save(parking).Error
}

func (r *ParkingRepository) CreateParkingImages(images []entities.ParkingImage) error {
	if len(images) == 0 {
		return nil
	}

	return r.db.Create(&images).Error
}

func (r *ParkingRepository) FindByIDAndUserID(id uint, userID uint) (*entities.Parking, error) {
	var parking entities.Parking
	// Preload کردن Images برای اینکه موقع حذف بتوانیم کلیدهای تصاویر را هم داشته باشیم
	err := r.db.Preload("Images").Where("id = ? AND user_id = ?", id, userID).First(&parking).Error
	if err != nil {
		return nil, err
	}
	return &parking, nil
}

func (r *ParkingRepository) Delete(parking *entities.Parking) error {
	// اگر از gorm.Model استفاده می‌کنی این دستور به صورت خودکار Soft Delete انجام می‌دهد
	return r.db.Delete(parking).Error
}
