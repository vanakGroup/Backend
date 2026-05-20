package service

import (
	"github.com/Vanak/Backend/bootstrap"
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/ports"
)

type ParkingService struct {
	unitOfWork ports.UnitOfWork
	cfg        *bootstrap.Config
	logger     logging.Logger
}

func NewParkingService(cfg *bootstrap.Config, unitOfWork ports.UnitOfWork) *ParkingService {
	logger := logging.NewLogger(cfg)

	return &ParkingService{
		logger:     logger,
		cfg:        cfg,
		unitOfWork: unitOfWork,
	}
}

func (ps *ParkingService) CreateParking(userID uint, parkingInfo parkingdto.CreateParkingRequest) (*parkingdto.CreateParkingResponse, error) {
	err := ps.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		parkingRepo := rf.ParkingRepository()
		addressRepo := rf.AddressRepository()

		parking := &entities.Parking{
			Title:     parkingInfo.Title,
			Latitude:  parkingInfo.Latitude,
			Longitude: parkingInfo.Longitude,
			BasePrice: parkingInfo.BasePrice,
			Status:    enums.ParkingStatusActive,
			UserID:    userID,
		}

		err := parkingRepo.CreateParking(parking)

		if err != nil {
			return err
		}

		city, err := addressRepo.GetCityByID(parkingInfo.Address.CityID)

		if err != nil {
			return err
		}

		if city == nil {
			notFoundError := exception.NotFoundError{Item: ps.cfg.Constant.ErrorField.City}
			return notFoundError
		}

		address := &entities.Address{
			CityID:       city.ID,
			PostalCode:   parkingInfo.Address.PostalCode,
			UnitNumber:   parkingInfo.Address.UnitNumber,
			ExactAddress: parkingInfo.Address.ExactAddress,
			ParkingID:    parking.ID,
		}

		err = addressRepo.CreateAddress(address)

		if err != nil {
			return err
		}

		if len(parkingInfo.AttributeIDs) > 0 {
			attributes, err := parkingRepo.GetAttributesByIDs(parkingInfo.AttributeIDs)

			if err != nil {
				return err
			}

			if len(attributes) == 0 {
				notFoundError := exception.NotFoundError{Item: ps.cfg.Constant.ErrorField.Attribute}
				return notFoundError
			}

			err = parkingRepo.AddAttributes(parking, attributes)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return &parkingdto.CreateParkingResponse{}, err
	}

	return &parkingdto.CreateParkingResponse{}, nil
}
