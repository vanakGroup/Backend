package usecase

import (
	reservationdto "github.com/Vanak/Backend/internal/application/dto/reservation"
)

type ReservationService interface {
	CreateReservation(reservationdto.ReservationDTORequest) error
	GetReservationListForParking(parkingIDHashed string, userID uint) ([]reservationdto.ReservationDTOResponse, error)
	GetReservationDetailForParking(reserveIDHashed string, parkingIDHashed string, userID uint) (*reservationdto.ReservationDTOResponse, error)
	GetReservationListForVehicle(vehicleIDHashed string, userID uint) ([]reservationdto.ReservationDTOResponse, error)
	GetReservationDetailForVehicle(reserveIDHashed string, vehicleIDhashed string, userID uint) (*reservationdto.ReservationDTOResponse, error)
	VerifySettlement(reserveIDHashed, parkingIDHashed string, userID uint) error
}
