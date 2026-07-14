package postgres

import (
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewReservationRepository(cfg *bootstrap.Config, db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{
		cfg: cfg,
		db:  db,
	}
}

func (rr *ReservationRepository) FindAndLockSlots(parkingID uint, requiredHours []uint, targetDate time.Time) ([]entities.AvailabilitySlot, error) {
	var slots []entities.AvailabilitySlot

	dateString := targetDate.Format("2006-01-02")

	err := rr.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&entities.AvailabilitySlot{}).
		Where("parking_id = ? AND is_booked = false", parkingID).
		Where("DATE(date) = ? AND hour IN ?", dateString, requiredHours).
		Find(&slots).Error

	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (rr *ReservationRepository) Create(reservation *entities.Reservation) error {
	return rr.db.Model(&entities.Reservation{}).Create(reservation).Error
}

func (rr *ReservationRepository) UpdateSlotsWithReservation(slotIDs []uint, reservationID uint) error {
	return rr.db.Model(&entities.AvailabilitySlot{}).
		Where("id IN ?", slotIDs).
		Updates(map[string]interface{}{
			"is_booked":      true,
			"reservation_id": reservationID,
		}).Error
}

func (rr *ReservationRepository) FindParkingReservation(parkingID uint) ([]entities.Reservation, error) {
	var reservations []entities.Reservation

	err := rr.db.Model(&entities.Reservation{}).
		Where("parking_id = ?", parkingID).
		Preload("Vehicle").
		Preload("User").
		Find(&reservations).Error

	if err != nil {
		return nil, err
	}

	return reservations, nil
}
func (rr *ReservationRepository) FindVehicleReservation(vehicleID uint) ([]entities.Reservation, error) {
	var reservations []entities.Reservation
	err := rr.db.Model(&entities.Reservation{}).
		Where("vehicle_id = ?", vehicleID).
		Preload("User").
		Preload("Parking.Address").
		Find(&reservations).Error

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (rr *ReservationRepository) FindReservationByIDAndParkingID(reserveID uint, parkingID uint) (*entities.Reservation, error) {
	var reservation *entities.Reservation
	err := rr.db.Model(&entities.Reservation{}).
		Where("id = ?", reserveID).
		Where("parking_id = ?", parkingID).First(&reservation).Error

	if err != nil {
		return nil, err
	}

	return reservation, nil
}
func (rr *ReservationRepository) FindReservationByIDAndVehicleID(reserveID uint, vehicleID uint) (*entities.Reservation, error) {
	var reservation *entities.Reservation
	err := rr.db.Model(&entities.Reservation{}).
		Where("id = ?", reserveID).
		Where("vehicle_id = ?", vehicleID).
		First(&reservation).Error

	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (rr *ReservationRepository) UpdateReservationStatus(reserveID uint, newStatus enums.ReservationStatus) error {
	err := rr.db.Model(&entities.Reservation{}).
		Where("id = ?", reserveID).
		Update("status", newStatus).Error

	if err != nil {
		return err
	}

	return nil
}
