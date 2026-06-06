package vehicledto

import (
	"mime/multipart"
)

type VehicleDTORequest struct {
	LicensePlate string
	CarModel     string
	CardFile     *multipart.FileHeader
	CardFileName string
	UserId       uint
}
