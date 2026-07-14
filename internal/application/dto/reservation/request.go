package reservationdto

type ReservationDTORequest struct {
	StartTime int64
	EndTime   int64
	ParkingID string
	UserID    uint
	VehicleID string
}
