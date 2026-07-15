package service

import (
	"context"
	"strings"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	vehicledto "github.com/Vanak/Backend/internal/application/dto/vehicle"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/domain/ports/repository/postgres"
	"github.com/Vanak/Backend/internal/domain/storage"
	"github.com/Vanak/Backend/internal/domain/utils"
	"github.com/Vanak/Backend/internal/infrastructure/hashed"
)

type VehicleService struct {
	unitOfWork  ports.UnitOfWork
	cfg         *bootstrap.Config
	S3Storage   storage.S3Storage
	vehicleRepo postgres.VehicleRepository
	hashID      *hashed.HashID
}

func NewVehicleService(
	unitOfWork ports.UnitOfWork,
	cfg *bootstrap.Config,
	S3Storage storage.S3Storage,
	vehicleRepo postgres.VehicleRepository,
	hashID *hashed.HashID,
) *VehicleService {
	return &VehicleService{
		cfg:         cfg,
		S3Storage:   S3Storage,
		vehicleRepo: vehicleRepo,
		unitOfWork:  unitOfWork,
		hashID:      hashID,
	}
}

func (vs *VehicleService) RegisterVehicle(ctx context.Context, vehicleInfo vehicledto.VehicleDTORequest) error {
	formatError := exception.NewFormatError()
	alreadyExistsError := exception.NewAlreadyExistsError()
	validationError := exception.NewValidationErrors()

	if strings.TrimSpace(vehicleInfo.CarModel) == "" {
		validationError.Add(vs.cfg.Constant.ErrorField.CarModel, vs.cfg.Constant.ErrorTag.Empty)
	}

	if vehicleInfo.CardFile.Size > 100*5000*5000 {
		validationError.Add(vs.cfg.Constant.ErrorField.RegistrationCardSize, vs.cfg.Constant.ErrorTag.TooLarge)
	}

	if validationError.HasErrors() {
		return validationError
	}

	LicensePlate, check := utils.CheckLicensePlate(vehicleInfo.LicensePlate)
	if !check {
		formatError.Add(vs.cfg.Constant.VehicleFields.LicensePlate, vs.cfg.Constant.ErrorTag.WrongFormat, vs.cfg.Constant.ErrorField.LicensePlate)
	}

	if formatError.HasErrors() {
		return formatError
	}

	if vs.unitOfWork.Factory().VehicleRepository().FindVehiclesByLicensePlate(*LicensePlate) {
		alreadyExistsError.Add(vs.cfg.Constant.VehicleFields.LicensePlate, vs.cfg.Constant.ErrorTag.AlreadyExist, vs.cfg.Constant.ErrorField.LicensePlate)
	}

	if alreadyExistsError.HasErrors() {
		return alreadyExistsError
	}

	err := vs.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		vehicleRepo := rf.VehicleRepository()

		registrationCardPath := vs.cfg.Constant.S3BucketPath.GenerateVehicleImagePath(vehicleInfo.CardFileName)
		err := vs.S3Storage.UploadFile(enums.RegistrationCardPic, registrationCardPath, vehicleInfo.CardFile)
		if err != nil {
			return err
		}

		newVehicle := entities.Vehicle{
			LicensePlate:        *LicensePlate,
			CarModel:            vehicleInfo.CarModel,
			RegistrationCardUrl: registrationCardPath,
			UserID:              vehicleInfo.UserId,
		}

		if err := vehicleRepo.Create(&newVehicle); err != nil {
			vs.S3Storage.DeleteObject(enums.RegistrationCardPic, registrationCardPath)
			return err
		}

		return nil
	})

	return err
}

func (vs *VehicleService) GetVehicle(vehicleHashedId string, userId uint) (*vehicledto.VehicleDTOResponse, error) {
	vehicleId, err := vs.hashID.Decode(vehicleHashedId)
	if err != nil {
		return nil, err
	}

	vehicle, err := vs.unitOfWork.Factory().VehicleRepository().FindUserVehicle(vehicleId, userId)
	if err != nil {
		not_found := exception.NewNotFoundError(vs.cfg.Constant.ErrorField.Vehicle)
		return nil, not_found
	}

	expirationTime := 15 * time.Minute
	registrationCardUrl, errUrl := vs.S3Storage.GetPresignedURL(enums.RegistrationCardPic, vehicle.RegistrationCardUrl, expirationTime)
	if errUrl != nil {
		return nil, errUrl
	}

	vehicleInformation := vehicledto.VehicleDTOResponse{
		LicensePlate:        vehicle.LicensePlate,
		CarModel:            vehicle.CarModel,
		RegistrationCardUrl: registrationCardUrl,
	}
	return &vehicleInformation, nil
}

func (vs *VehicleService) DeleteVehicle(vehicleHashedId string, userId uint) error {
	vehicleId, err := vs.hashID.Decode(vehicleHashedId)
	if err != nil {
		return err
	}

	reservationList, verr := vs.unitOfWork.Factory().ReservationRepository().FindVehicleReservation(vehicleId)
	if verr != nil {
		return verr
	}

	if len(reservationList) >= 1 {
		return exception.NewConflictError(vs.cfg.Constant.ErrorField.Reserve)
	}

	vehicle, err := vs.vehicleRepo.FindUserVehicle(vehicleId, userId)
	if err != nil {
		not_found := exception.NewNotFoundError(vs.cfg.Constant.ErrorField.Vehicle)
		return not_found
	}

	err = vs.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		err := rf.VehicleRepository().RemoveVehicle(vehicleId, userId)
		if err != nil {
			return err
		}

		err = vs.S3Storage.DeleteObject(enums.RegistrationCardPic, vehicle.RegistrationCardUrl)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (vs *VehicleService) GetAllVehicles(userId uint) ([]vehicledto.VehicleListDTOResponse, error) {
	vehicles, err := vs.vehicleRepo.FindAllUserVehicles(userId)
	if err != nil {
		return nil, err
	}

	vehiclesList := make([]vehicledto.VehicleListDTOResponse, 0, len(vehicles))
	for _, v := range vehicles {
		vehicleID, err := vs.hashID.Encode(v.ID)
		if err != nil {
			return nil, err
		}

		expirationTime := 15 * time.Minute
		registrationCardUrl, errUrl := vs.S3Storage.GetPresignedURL(enums.RegistrationCardPic, v.RegistrationCardUrl, expirationTime)
		if errUrl != nil {
			return nil, errUrl
		}

		vehiclesList = append(vehiclesList, vehicledto.VehicleListDTOResponse{
			ID:                  vehicleID,
			LicensePlate:        v.LicensePlate,
			CarModel:            v.CarModel,
			RegistrationCardUrl: registrationCardUrl,
		})
	}

	return vehiclesList, nil
}
