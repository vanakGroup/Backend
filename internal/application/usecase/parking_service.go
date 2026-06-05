package usecase

import parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"

type ParkingService interface {
	CreateParking(userID uint, parkingInfo parkingdto.CreateParkingRequest) (*parkingdto.CreateParkingResponse, error)
	DeleteParking(userID uint, parkingID uint) error
	GetUserParkings(userID uint) ([]parkingdto.ParkingResponse, error)
	UpdateParking(userID uint, parkingInfo parkingdto.UpdateParkingRequest) error
}
