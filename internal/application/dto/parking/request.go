package parkingdto

import (
	"mime/multipart"
	"time"

	addressdto "github.com/Vanak/Backend/internal/application/dto/address"
)

type CreateParkingRequest struct {
	Title     string
	Latitude  float64
	Longitude float64
	BasePrice float64
	Address   addressdto.AddressRequest

	AttributeIDs []uint
	Slots        []SlotRequest

	DocumentFile *multipart.FileHeader
	Images       []*multipart.FileHeader
}

type SlotRequest struct {
	Date      time.Time
	StartHour uint
	EndHour   uint
}

type UpdateParkingRequest struct {
	ParkingID string
	Title     string
	Latitude  float64
	Longitude float64
	BasePrice float64
	Address   addressdto.UpdateAddressRequest

	AttributeIDs []uint
	Slots        []SlotRequest

	DocumentFile *multipart.FileHeader
	Images       []*multipart.FileHeader
}

type GetParkingByIDAndUserIDRequest struct {
	UserID          uint
	ParkingIDHashed string
}
