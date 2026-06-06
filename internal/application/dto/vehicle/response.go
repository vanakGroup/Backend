package vehicledto

type VehicleDTOResponse struct {
	LicensePlate        string `json:"license_plate"`
	CarModel            string `json:"car_model"`
	RegistrationCardUrl string `json:"registration_card_url"`
}

type VehicleListDTOResponse struct {
	ID           uint   `json:"id"`
	LicensePlate string `json:"license_plate"`
	CarModel     string `json:"car_model"`
}
