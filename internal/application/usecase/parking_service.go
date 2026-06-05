package usecase

import parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"

type ParkingService interface {
	CreateParking(userID uint, parkingInfo parkingdto.CreateParkingRequest) (*parkingdto.CreateParkingResponse, error)
	DeleteParking(userID uint, parkingID uint) error
	GetParkingDetail(userID uint, parkingID uint) (*parkingdto.ParkingDetailResponse, error)
	UpdateParking(userID uint, parkingInfo parkingdto.UpdateParkingRequest) error
}
