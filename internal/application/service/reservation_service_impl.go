package service

import (
	"time"

	"github.com/Vanak/Backend/bootstrap"
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	reservationdto "github.com/Vanak/Backend/internal/application/dto/reservation"
	vehicledto "github.com/Vanak/Backend/internal/application/dto/vehicle"
	walletdto "github.com/Vanak/Backend/internal/application/dto/wallet"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/infrastructure/hashed"
)

type ReservationService struct {
	cfg            *bootstrap.Config
	unitOfWork     ports.UnitOfWork
	hashID         *hashed.HashID
	parkingService usecase.ParkingService
	vehicleService usecase.VehicleService
	walletService  usecase.WalletService
}

func NewReservationService(
	cfg *bootstrap.Config,
	unitOfWork ports.UnitOfWork,
	hashID *hashed.HashID,
	parkingService usecase.ParkingService,
	vehicleService usecase.VehicleService,
	walletService usecase.WalletService,

) *ReservationService {
	return &ReservationService{
		cfg:            cfg,
		unitOfWork:     unitOfWork,
		hashID:         hashID,
		parkingService: parkingService,
		vehicleService: vehicleService,
		walletService:  walletService,
	}
}

func (rs *ReservationService) CreateReservation(reservationInfo reservationdto.ReservationDTORequest) error {
	userID := reservationInfo.UserID
	disabledErr := exception.NewDisabledError()

	wallet, err := rs.unitOfWork.Factory().WalletRepository().FindWalletByUserId(userID)
	if err != nil {
		return err
	}

	if wallet.IsActive == false {
		disabledErr.Add(rs.cfg.Constant.ErrorField.Wallet, rs.cfg.Constant.ErrorTag.Disabled)
	}

	parkingID, err := rs.hashID.Decode(reservationInfo.ParkingID)
	if err != nil {
		return err
	}

	vehicleID, err := rs.hashID.Decode(reservationInfo.VehicleID)
	if err != nil {
		return err
	}

	_, err = rs.unitOfWork.Factory().VehicleRepository().FindUserVehicle(vehicleID, userID)
	if err != nil {
		return err
	}

	parking, err := rs.unitOfWork.Factory().ParkingRepository().FindByID(parkingID)
	if err != nil {
		return err
	}

	if parking.Status == enums.ParkingStatusInActive {
		disabledErr.Add(rs.cfg.Constant.ErrorField.Parking, rs.cfg.Constant.ErrorTag.Disabled)
	}

	if disabledErr.HasErrors() {
		return disabledErr
	}

	vehicleReservations, err := rs.unitOfWork.Factory().ReservationRepository().FindVehicleReservation(vehicleID)
	if err != nil {
		return err
	}

	startTimeInt := reservationInfo.StartTime
	endTimeInt := reservationInfo.EndTime

	for _, r := range vehicleReservations {
		if startTimeInt < r.EndTime && endTimeInt > r.StartTime {
			return exception.NewConflictError(rs.cfg.Constant.ErrorField.Reserve)
		}
	}

	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return err
	}

	startT := time.Unix(startTimeInt, 0).In(loc)
	endT := time.Unix(endTimeInt, 0).In(loc)

	pureDate := time.Date(startT.Year(), startT.Month(), startT.Day(), 0, 0, 0, 0, loc)

	var requiredHours []uint
	for t := startT; t.Before(endT); t = t.Add(time.Hour) {
		requiredHours = append(requiredHours, uint(t.Hour()))
	}

	err = rs.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		reservationRepo := rf.ReservationRepository()
		slots, err := reservationRepo.FindAndLockSlots(parkingID, requiredHours, pureDate)
		if err != nil {
			return err
		}

		if len(slots) != len(requiredHours) {
			return &exception.ConflictError{
				Item: rs.cfg.Constant.ErrorField.Reserve,
			}
		}

		newReservation := &entities.Reservation{
			TotalPrice: int64(parking.BasePrice) * int64(len(slots)),
			StartTime:  reservationInfo.StartTime,
			EndTime:    reservationInfo.EndTime,
			Hours:      len(requiredHours),
			Status:     enums.CONFIRMED,
			UserID:     reservationInfo.UserID,
			ParkingID:  parkingID,
			VehicleID:  vehicleID,
			BasePrice:  int64(parking.BasePrice),
		}

		err = reservationRepo.Create(newReservation)
		if err != nil {
			return err
		}

		slotIDs := make([]uint, 0)
		for _, s := range slots {
			slotIDs = append(slotIDs, s.ID)
		}

		err = reservationRepo.UpdateSlotsWithReservation(slotIDs, newReservation.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (rs *ReservationService) GetReservationListForParking(parkingIDHashed string, userID uint) ([]reservationdto.ReservationDTOResponse, error) {

	parkingID, err := rs.hashID.Decode(parkingIDHashed)

	if err != nil {
		return nil, err
	}

	parking, err := rs.parkingService.GetParkingDetail(userID, parkingIDHashed)
	if err != nil {
		return nil, err
	}

	reservationList, err := rs.unitOfWork.Factory().ReservationRepository().FindParkingReservation(parkingID)
	if err != nil {
		return nil, err
	}

	reservationResponseList := make([]reservationdto.ReservationDTOResponse, 0)
	for _, r := range reservationList {
		id, err := rs.hashID.Encode(r.ID)
		if err != nil {
			return nil, err
		}

		vehicle := r.Vehicle

		user := r.User

		vehicleIdHashed, err := rs.hashID.Encode(r.VehicleID)
		if err != nil {
			return nil, err
		}

		reservationResponseList = append(reservationResponseList, reservationdto.ReservationDTOResponse{
			ID:         id,
			FirstName:  *user.FirstName,
			LastName:   *user.LastName,
			TotalPrice: r.TotalPrice,
			Baseprice:  r.BasePrice,
			StartTime:  r.StartTime,
			EndTime:    r.EndTime,
			Hours:      r.Hours,
			Status:     r.Status.String(),
			Parking: parkingdto.ParkingReservationDetailResponse{
				ID:        parkingIDHashed,
				Title:     parking.Title,
				Latitude:  parking.Latitude,
				Longitude: parking.Longitude,
				Address:   parking.Address,
			},
			Vehicle: vehicledto.VehicleListDTOResponse{
				ID:           vehicleIdHashed,
				CarModel:     vehicle.CarModel,
				LicensePlate: vehicle.LicensePlate,
			},
		})
	}

	return reservationResponseList, nil
}

func (rs *ReservationService) GetReservationDetailForParking(reserveIDHashed string, parkingIDHashed string, userID uint) (*reservationdto.ReservationDTOResponse, error) {
	reserveID, err := rs.hashID.Decode(reserveIDHashed)
	if err != nil {
		return nil, err
	}

	parkingID, err := rs.hashID.Decode(parkingIDHashed)
	if err != nil {
		return nil, err
	}

	reservation, err := rs.unitOfWork.Factory().ReservationRepository().FindReservationByIDAndParkingID(reserveID, parkingID)
	if err != nil {
		return nil, err
	}

	parking, err := rs.parkingService.GetParkingDetail(userID, parkingIDHashed)
	if err != nil {
		return nil, err
	}

	user, err := rs.unitOfWork.Factory().UserRepository().FindUserById(reservation.UserID)
	if err != nil {
		return nil, err
	}

	vehicle, err := rs.unitOfWork.Factory().VehicleRepository().FindVehicleByID(reservation.VehicleID)
	if err != nil {
		return nil, err
	}

	vehicleIDHashed, err := rs.hashID.Encode(reservation.VehicleID)
	if err != nil {
		return nil, err
	}

	reservationReponse := reservationdto.ReservationDTOResponse{
		ID:         reserveIDHashed,
		FirstName:  *user.FirstName,
		LastName:   *user.LastName,
		TotalPrice: reservation.TotalPrice,
		StartTime:  reservation.StartTime,
		EndTime:    reservation.EndTime,
		Hours:      reservation.Hours,
		Status:     reservation.Status.String(),
		Parking: parkingdto.ParkingReservationDetailResponse{
			ID:        parkingIDHashed,
			Title:     parking.Title,
			Latitude:  parking.Latitude,
			Longitude: parking.Longitude,
			Address:   parking.Address,
		},
		Vehicle: vehicledto.VehicleListDTOResponse{
			ID:           vehicleIDHashed,
			CarModel:     vehicle.CarModel,
			LicensePlate: vehicle.LicensePlate,
		},
	}

	return &reservationReponse, nil
}

func (rs *ReservationService) GetReservationListForVehicle(vehicleIDHashed string, userID uint) ([]reservationdto.ReservationDTOResponse, error) {
	vehicleID, err := rs.hashID.Decode(vehicleIDHashed)
	if err != nil {
		return nil, err
	}

	vehicle, err := rs.unitOfWork.Factory().VehicleRepository().FindUserVehicle(vehicleID, userID)
	if err != nil {
		return nil, err
	}

	reservationList, err := rs.unitOfWork.Factory().ReservationRepository().FindVehicleReservation(vehicleID)
	if err != nil {
		return nil, err
	}

	reservationResponseList := make([]reservationdto.ReservationDTOResponse, 0)
	for _, r := range reservationList {
		id, err := rs.hashID.Encode(r.ID)
		if err != nil {
			return nil, err
		}

		user := r.User

		parkingIDHashed, err := rs.hashID.Encode(r.ParkingID)
		if err != nil {
			return nil, err
		}

		parking := r.Parking

		reservationResponseList = append(reservationResponseList, reservationdto.ReservationDTOResponse{
			ID:         id,
			FirstName:  *user.FirstName,
			LastName:   *user.LastName,
			TotalPrice: r.TotalPrice,
			Baseprice:  r.BasePrice,
			StartTime:  r.StartTime,
			EndTime:    r.EndTime,
			Hours:      r.Hours,
			Status:     r.Status.String(),
			Parking: parkingdto.ParkingReservationDetailResponse{
				ID:        parkingIDHashed,
				Title:     parking.Title,
				Latitude:  parking.Latitude,
				Longitude: parking.Longitude,
				Address: parkingdto.AddressDetailResponse{
					CityID:       parking.Address.CityID,
					PostalCode:   parking.Address.PostalCode,
					UnitNumber:   parking.Address.UnitNumber,
					ExactAddress: parking.Address.ExactAddress,
				},
			},
			Vehicle: vehicledto.VehicleListDTOResponse{
				ID:           vehicleIDHashed,
				CarModel:     vehicle.CarModel,
				LicensePlate: vehicle.LicensePlate,
			},
		})
	}

	return reservationResponseList, nil
}

func (rs *ReservationService) GetReservationDetailForVehicle(reserveIDHashed string, vehicleIDHashed string, userID uint) (*reservationdto.ReservationDTOResponse, error) {
	reserveID, err := rs.hashID.Decode(reserveIDHashed)
	if err != nil {
		return nil, err
	}

	vehicleID, err := rs.hashID.Decode(vehicleIDHashed)
	if err != nil {
		return nil, err
	}

	vehicle, err := rs.unitOfWork.Factory().VehicleRepository().FindUserVehicle(vehicleID, userID)
	if err != nil {
		return nil, err
	}

	reservation, err := rs.unitOfWork.Factory().ReservationRepository().FindReservationByIDAndVehicleID(reserveID, vehicleID)
	if err != nil {
		return nil, err
	}

	user, err := rs.unitOfWork.Factory().UserRepository().FindUserById(userID)
	if err != nil {
		return nil, err
	}

	parkingEntity, err := rs.unitOfWork.Factory().ParkingRepository().FindByID(reservation.ParkingID)
	if err != nil {
		return nil, err
	}

	parkingIDHashed, err := rs.hashID.Encode(reservation.ParkingID)
	if err != nil {
		return nil, err
	}

	parking, err := rs.parkingService.GetParkingDetail(parkingEntity.UserID, parkingIDHashed)
	if err != nil {
		return nil, err
	}

	reservationReponse := reservationdto.ReservationDTOResponse{
		ID:         reserveIDHashed,
		FirstName:  *user.FirstName,
		LastName:   *user.LastName,
		TotalPrice: reservation.TotalPrice,
		StartTime:  reservation.StartTime,
		EndTime:    reservation.EndTime,
		Hours:      reservation.Hours,
		Status:     reservation.Status.String(),
		Parking: parkingdto.ParkingReservationDetailResponse{
			ID:        parkingIDHashed,
			Title:     parking.Title,
			Latitude:  parking.Latitude,
			Longitude: parking.Longitude,
			Address:   parking.Address,
		},
		Vehicle: vehicledto.VehicleListDTOResponse{
			ID:           vehicleIDHashed,
			CarModel:     vehicle.CarModel,
			LicensePlate: vehicle.LicensePlate,
		},
	}

	return &reservationReponse, nil

}
func (rs *ReservationService) VerifySettlement(reserveIDHashed, parkingIDHashed string, userID uint) error {
	reserveID, err := rs.hashID.Decode(reserveIDHashed)
	if err != nil {
		return err
	}

	parkingID, err := rs.hashID.Decode(parkingIDHashed)
	if err != nil {
		return err
	}

	parking, err := rs.unitOfWork.Factory().ParkingRepository().FindByIDAndUserID(parkingID, userID)
	if err != nil {
		return err
	}

	reservation, err := rs.unitOfWork.Factory().ReservationRepository().FindReservationByIDAndParkingID(reserveID, parkingID)

	if err != nil {
		return err
	}

	currentUnixTime := time.Now().Unix()
	if currentUnixTime < reservation.EndTime {
		return exception.NewConflictError(rs.cfg.Constant.ErrorField.Time)
	}

	deductDTO := walletdto.WalletDTORequest{
		Amount:          reservation.TotalPrice,
		ReferenceType:   enums.RESERVATION,
		TransactionType: enums.PAYMENT,
		UserId:          reservation.UserID,
		ReferenceId:     &reservation.ID,
	}

	depoistDTO := walletdto.WalletDTORequest{
		Amount:          reservation.TotalPrice,
		ReferenceType:   enums.RESERVATION,
		TransactionType: enums.DEPOSIT,
		UserId:          parking.UserID,
		ReferenceId:     &reservation.ID,
	}

	err = rs.walletService.Transfer(deductDTO, depoistDTO)
	if err == nil {
		rs.unitOfWork.Factory().ReservationRepository().UpdateReservationStatus(reserveID, enums.COMPLETED)
	}

	return err
}
