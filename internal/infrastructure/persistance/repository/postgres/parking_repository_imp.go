package postgres

import (
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/filter"
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

func (pr *ParkingRepository) GetAttributeByName(name string) (*entities.ParkingAttribute, error) {
	var attribute entities.ParkingAttribute
	result := pr.db.Where("name = ?", name).First(&attribute)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &attribute, nil
}

func (pr *ParkingRepository) CreateAttribute(parkingAttribute *entities.ParkingAttribute) error {
	return pr.db.Create(parkingAttribute).Error
}

func (pr *ParkingRepository) CreateParking(parking *entities.Parking) error {
	return pr.db.Create(parking).Error
}

func (pr *ParkingRepository) GetAttributesByIDs(ids []uint) ([]*entities.ParkingAttribute, error) {
	var attributes []*entities.ParkingAttribute
	result := pr.db.Where("id IN ?", ids).Find(&attributes)
	if result.Error != nil {
		return nil, result.Error
	}
	return attributes, nil
}

func (pr *ParkingRepository) AddAttributes(parking *entities.Parking, attributes []*entities.ParkingAttribute) error {
	return pr.db.Model(parking).Association("Attributes").Append(attributes)
}

func (pr *ParkingRepository) CreateSlots(slots []entities.AvailabilitySlot) error {
	return pr.db.Create(&slots).Error
}

func (pr *ParkingRepository) UpdateParking(parking *entities.Parking) error {
	return pr.db.Save(parking).Error
}

func (pr *ParkingRepository) CreateParkingImages(images []entities.ParkingImage) error {
	if len(images) == 0 {
		return nil
	}

	return pr.db.Create(&images).Error
}

func (pr *ParkingRepository) FindByIDAndUserID(id uint, userID uint) (*entities.Parking, error) {
	var parking entities.Parking
	err := pr.db.Preload("Images").Where("id = ? AND user_id = ?", id, userID).First(&parking).Error
	if err != nil {
		return nil, err
	}
	return &parking, nil
}

func (pr *ParkingRepository) Delete(parking *entities.Parking) error {
	return pr.db.Transaction(func(tx *gorm.DB) error {
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

func (pr *ParkingRepository) GetSlotsByParkingID(parkingID uint) ([]entities.AvailabilitySlot, error) {
	var slots []entities.AvailabilitySlot
	err := pr.db.Where("parking_id = ?", parkingID).Find(&slots).Error
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (pr *ParkingRepository) DeleteSlotsWithIDs(slotIDs []uint) error {
	if len(slotIDs) == 0 {
		return nil
	}

	return pr.db.Where("id IN ?", slotIDs).Delete(&entities.AvailabilitySlot{}).Error
}

func (pr *ParkingRepository) DeleteParkingImagesWithIDs(imageIDs []uint) error {
	if len(imageIDs) == 0 {
		return nil
	}
	return pr.db.Where("id IN ?", imageIDs).Delete(&entities.ParkingImage{}).Error
}

func (pr *ParkingRepository) UpdateAttributes(parking *entities.Parking, attributes []entities.ParkingAttribute) error {
	return pr.db.Model(parking).Association("Attributes").Replace(attributes)
}

func (pr *ParkingRepository) GetDetailByIDAndUserID(parkingID uint, userID uint) (*entities.Parking, error) {
	var parking entities.Parking

	err := pr.db.Where("id = ? AND user_id = ?", parkingID, userID).
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

func (pr *ParkingRepository) GetAllAttributes() ([]entities.ParkingAttribute, error) {
	var attributes []entities.ParkingAttribute
	err := pr.db.Select("id, name").Find(&attributes).Error
	return attributes, err
}

func (pr *ParkingRepository) GetParkingByIDAndUserID(userID uint, parkingID uint) (*entities.Parking, error) {
	var parking entities.Parking

	err := pr.db.Preload("Address").
		Preload("Images").
		Preload("Attributes").
		Where("id = ? AND user_id = ?", parkingID, userID).
		First(&parking).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &parking, nil
}

func (pr *ParkingRepository) GetAllByUserID(userID uint) ([]*entities.Parking, error) {
	var parkings []*entities.Parking

	err := pr.db.Preload("Address").
		Preload("Images").
		Preload("Attributes").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&parkings).Error

	if err != nil {
		return nil, err
	}

	return parkings, nil
}

func (pr *ParkingRepository) GetAll(options *filter.QueryOptions) ([]entities.Parking, error) {
	var parkings []entities.Parking
	query := pr.db.Model(&entities.Parking{}).
		Preload("Address").
		Preload("Images").
		Preload("Attributes")

	query = applyQueryOptions(query, options)

	err := query.Find(&parkings).Error
	if err != nil {
		return nil, err
	}
	return parkings, nil
}

func (pr *ParkingRepository) FindByID(id uint) (*entities.Parking, error) {
	var parking entities.Parking
	err := pr.db.Model(&entities.Parking{}).
		Where("id = ?", id).First(&parking).Error

	if err != nil {
		return nil, err
	}

	return &parking, err
}

func (pr *ParkingRepository) CountAll() (int64, error) {
	var total int64
	err := pr.db.Model(&entities.Parking{}).Count(&total).Error
	return total, err
}

func (pr *ParkingRepository) GetByID(id uint) (*entities.Parking, error) {
	var parking entities.Parking
	err := pr.db.Preload("Address").
		Preload("Images").
		Preload("Attributes").
		First(&parking, id).Error
	if err != nil {
		return nil, err
	}
	return &parking, nil
}

func (pr *ParkingRepository) Count(f filter.ParkingFilter) (int64, error) {
	baseQuery := pr.db.Model(&entities.Parking{})

	query := NewParkingQueryBuilder(baseQuery).
		Apply(f).
		Build()

	var total int64
	err := query.Count(&total).Error
	return total, err
}

func (pr *ParkingRepository) Search(f filter.ParkingFilter, options *filter.QueryOptions) ([]entities.Parking, error) {
	baseQuery := pr.db.Model(&entities.Parking{})

	query := NewParkingQueryBuilder(baseQuery).
		Apply(f).
		Build()

	var parkings []entities.Parking

	query = applyQueryOptions(query, options)

	err := query.
		Preload("Address").
		Preload("Attributes").
		Preload("Images").
		Find(&parkings).Error

	if err != nil {
		return nil, err
	}

	return parkings, nil

}
