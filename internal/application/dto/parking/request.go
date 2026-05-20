package parkingdto

import (
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
}

type SlotRequest struct {
	StartTime time.Time
	EndTime   time.Time
}
