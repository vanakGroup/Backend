package usecase

import (
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	staticdto "github.com/Vanak/Backend/internal/application/dto/static"
)

type ParkingService interface {
	CreateParking(userID uint, parkingInfo parkingdto.CreateParkingRequest) (*parkingdto.CreateParkingResponse, error)
	DeleteParking(userID uint, parkingIDHashed string) error
	GetParkingDetail(userID uint, parkingIDHashed string) (*parkingdto.ParkingDetailResponse, error)
	UpdateParking(userID uint, parkingInfo parkingdto.UpdateParkingRequest) error
	GetCitiesList() ([]staticdto.StaticItemResponse, error)
	GetAttributesList() ([]staticdto.StaticItemResponse, error)
	GetParkingByIDAndUserID(req parkingdto.GetParkingByIDAndUserIDRequest) (*parkingdto.ParkingInfoResponse, error)
	GetUserParkings(userID uint) ([]parkingdto.ParkingForListResponse, error)
	SearchParking(req parkingdto.ParkingSearchRequest) ([]parkingdto.ParkingForListResponse, int64, error)
	GetAllParkings(limit, offset int) ([]parkingdto.ParkingForListResponse, int64, error)
	GetParkingPublicInfo(parkingIDHashed string) (*parkingdto.ParkingInfoResponse, error)
}
