package postgres

import (
	"time"

	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
)

type ReservationRepository interface {
	Create(*entities.Reservation) error
	FindAndLockSlots(parkingID uint, requiredHours []uint, date time.Time) ([]entities.AvailabilitySlot, error)
	UpdateSlotsWithReservation(slotIDs []uint, reservationID uint) error
	FindParkingReservation(parkingID uint) ([]entities.Reservation, error)
	FindVehicleReservation(vehicleID uint) ([]entities.Reservation, error)
	FindReservationByIDAndParkingID(reserveID uint, parkingID uint) (*entities.Reservation, error)
	FindReservationByIDAndVehicleID(reserveID uint, vehicleID uint) (*entities.Reservation, error)
	UpdateReservationStatus(reserveID uint, newStatus enums.ReservationStatus) error
}
