package postgres

import (
	"time"

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
	err := r.db.Preload("Images").Where("id = ? AND user_id = ?", id, userID).First(&parking).Error
	if err != nil {
		return nil, err
	}
	return &parking, nil
}

func (r *ParkingRepository) Delete(parking *entities.Parking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("parking_id = ?", parking.ID).Delete(&entities.AvailabilitySlot{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(parking).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *ParkingRepository) GetByUserID(userID uint) ([]entities.Parking, error) {
	var parkings []entities.Parking

	err := r.db.Preload("Address").Preload("Images").Where("user_id = ?", userID).Find(&parkings).Error
	if err != nil {
		return nil, err
	}

	return parkings, nil
}

func (r *ParkingRepository) HasActiveReservations(parkingID uint) (bool, error) {
	var count int64
	now := time.Now()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	currentHour := uint(now.Hour())

	err := r.db.Model(&entities.AvailabilitySlot{}).
		Where("parking_id = ? AND is_booked = ? AND (date > ? OR (date = ? AND hour >= ?))",
			parkingID,
			true,
			today,
			today,
			currentHour,
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ParkingRepository) GetSlotsByParkingID(parkingID uint) ([]entities.AvailabilitySlot, error) {
	var slots []entities.AvailabilitySlot
	err := r.db.Where("parking_id = ?", parkingID).Find(&slots).Error
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *ParkingRepository) DeleteSlotsWithIDs(slotIDs []uint) error {
	if len(slotIDs) == 0 {
		return nil
	}

	return r.db.Where("id IN ?", slotIDs).Delete(&entities.AvailabilitySlot{}).Error
}

func (r *ParkingRepository) DeleteParkingImagesWithIDs(imageIDs []uint) error {
	if len(imageIDs) == 0 {
		return nil
	}
	return r.db.Where("id IN ?", imageIDs).Delete(&entities.ParkingImage{}).Error
}

func (r *ParkingRepository) UpdateAttributes(parking *entities.Parking, attributes []entities.ParkingAttribute) error {
	return r.db.Model(parking).Association("Attributes").Replace(attributes)
}

func (r *ParkingRepository) GetDetailByIDAndUserID(parkingID uint, userID uint) (*entities.Parking, error) {
	var parking entities.Parking

	err := r.db.Where("id = ? AND user_id = ?", parkingID, userID).
		Preload("Address").
		Preload("Images").
		Preload("Attributes").
		Preload("Slots").
		First(&parking).Error

	if err != nil {
		return nil, err
	}
	return &parking, nil
}

func (r *ParkingRepository) GetAllAttributes() ([]entities.ParkingAttribute, error) {
	var attributes []entities.ParkingAttribute
	err := r.db.Select("id, name").Find(&attributes).Error
	return attributes, err
}