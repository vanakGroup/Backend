package reservationdto

import (
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	vehicledto "github.com/Vanak/Backend/internal/application/dto/vehicle"
)

type ReservationDTOResponse struct {
	ID         string                                      `json:"id"`
	FirstName  string                                      `json:"first_name"`
	LastName   string                                      `json:"last_name"`
	TotalPrice int64                                       `json:"total_price"`
	StartTime  int64                                       `json:"start_time"`
	EndTime    int64                                       `json:"end_time"`
	Hours      int                                         `json:"hours"`
	Status     string                                      `json:"status"`
	Parking    parkingdto.ParkingReservationDetailResponse `json:"parking"`
	Vehicle    vehicledto.VehicleListDTOResponse           `json:"vehicle"`
	Baseprice  int64                                       `json:"base_price"`
}
