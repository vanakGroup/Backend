package service

import (
	"fmt"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/domain/storage"
)

type ParkingService struct {
	unitOfWork ports.UnitOfWork
	cfg        *bootstrap.Config
	logger     logging.Logger
	s3Storage  storage.S3Storage
}

func NewParkingService(cfg *bootstrap.Config, unitOfWork ports.UnitOfWork, s3Storage storage.S3Storage) *ParkingService {
	logger := logging.NewLogger(cfg)

	return &ParkingService{
		logger:     logger,
		cfg:        cfg,
		unitOfWork: unitOfWork,
		s3Storage:  s3Storage,
	}
}

func (ps *ParkingService) CreateParking(userID uint, parkingInfo parkingdto.CreateParkingRequest) (*parkingdto.CreateParkingResponse, error) {
	var response *parkingdto.CreateParkingResponse
	var slotCount int

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
			return exception.NotFoundError{Item: ps.cfg.Constant.ErrorField.City}
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
				return exception.NotFoundError{Item: ps.cfg.Constant.ErrorField.Attribute}
			}

			err = parkingRepo.AddAttributes(parking, attributes)
			if err != nil {
				return err
			}
		}

		if len(parkingInfo.Slots) > 0 {
			allSlots, err := ps.createFirstSlotsForParking(parking.ID, parkingInfo.Slots)
			if err != nil {
				return err
			}

			if err := parkingRepo.CreateSlots(allSlots); err != nil {
				return err
			}

			slotCount = len(allSlots)
		}

		docKey := ps.cfg.Constant.S3BucketPath.GetParkingDocPath(parking.ID, parkingInfo.DocumentFile.Filename)
		err = ps.s3Storage.UploadFile(enums.ParkingDoc, docKey, parkingInfo.DocumentFile)
		if err != nil {
			return err
		}
		parking.DocumentURL = docKey

		var parkingImages []entities.ParkingImage
		for idx, fileHeader := range parkingInfo.Images {
			if fileHeader == nil {
				continue
			}

			imageKey := ps.cfg.Constant.S3BucketPath.GetParkingImagePath(parking.ID, idx, fileHeader.Filename)

			err = ps.s3Storage.UploadFile(enums.ParkingPic, imageKey, fileHeader)
			if err != nil {
				return err
			}

			parkingImages = append(parkingImages, entities.ParkingImage{
				ImageURL:  imageKey,
				ParkingID: parking.ID,
			})
		}

		err = parkingRepo.UpdateParking(parking)
		if err != nil {
			return err
		}

		if len(parkingImages) > 0 {
			err = parkingRepo.CreateParkingImages(parkingImages)
			if err != nil {
				return err
			}
		}

		var imgKeys []string
		for _, img := range parkingImages {
			imgKeys = append(imgKeys, img.ImageURL)
		}

		response = &parkingdto.CreateParkingResponse{
			ID:          parking.ID,
			Title:       parking.Title,
			Status:      parking.Status.String(),
			CreatedAt:   parking.CreatedAt,
			AddressID:   address.ID,
			SlotCount:   slotCount,
			DocumentURL: parking.DocumentURL,
			ImageURLs:   imgKeys,
		}

		return nil
	})

	if err != nil {
		return &parkingdto.CreateParkingResponse{}, err
	}

	return response, nil
}

func (ps *ParkingService) createFirstSlotsForParking(id uint, slotRequests []parkingdto.SlotRequest) ([]entities.AvailabilitySlot, error) {
	var allSlots []entities.AvailabilitySlot
	existingMap := make(map[string]bool)

	for _, slot := range slotRequests {
		for hour := slot.StartHour; hour < slot.EndHour; hour++ {
			key := fmt.Sprintf("%s:%d", slot.Date.Format("2006-01-02"), hour)
			if existingMap[key] {
				return nil, &exception.DuplicateError{Item: ps.cfg.Constant.ErrorField.Slot}
			}
			existingMap[key] = true

			newSlot := entities.AvailabilitySlot{
				Date:      slot.Date,
				Hour:      hour,
				IsBooked:  false,
				ParkingID: id,
			}

			allSlots = append(allSlots, newSlot)
		}
	}

	return allSlots, nil
}

func (ps *ParkingService) DeleteParking(userID uint, parkingID uint) error {
	parking, err := ps.unitOfWork.Factory().ParkingRepository().FindByIDAndUserID(parkingID, userID)
	if err != nil {
		return err
	}

	hasReservations, err := ps.unitOfWork.Factory().ParkingRepository().HasActiveReservations(parkingID)
	if err != nil {
		return err
	}
	if hasReservations {
		return &exception.ConflictError{
			Item: ps.cfg.Constant.ErrorField.Reserve,
		}
	}

	if parking.DocumentURL != "" {
		err := ps.s3Storage.DeleteObject(enums.ParkingDoc, parking.DocumentURL)
		if err != nil {
			return err
		}
	}

	if len(parking.Images) > 0 {
		for _, img := range parking.Images {
			if img.ImageURL != "" {
				err := ps.s3Storage.DeleteObject(enums.ParkingPic, img.ImageURL)
				if err != nil {
					return err
				}
			}
		}
	}

	return ps.unitOfWork.Factory().ParkingRepository().Delete(parking)
}

func (ps *ParkingService) GetUserParkings(userID uint) ([]parkingdto.ParkingResponse, error) {
	parkings, err := ps.unitOfWork.Factory().ParkingRepository().GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []parkingdto.ParkingResponse

	for _, p := range parkings {
		docURL := ""
		if p.DocumentURL != "" {
			url, err := ps.s3Storage.GetPresignedURL(enums.ParkingDoc, p.DocumentURL, time.Hour)
			if err == nil {
				docURL = url
			}
		}

		var imageUrls []string
		for _, img := range p.Images {
			if img.ImageURL != "" {
				url, err := ps.s3Storage.GetPresignedURL(enums.ParkingPic, img.ImageURL, time.Hour)
				if err == nil {
					imageUrls = append(imageUrls, url)
				}
			}
		}

		parkingDTO := parkingdto.ParkingResponse{
			ID:           p.ID,
			Title:        p.Title,
			Latitude:     p.Latitude,
			Longitude:    p.Longitude,
			BasePrice:    uint(p.BasePrice),
			DocumentURL:  docURL,
			ImageURLs:    imageUrls,
			ExactAddress: p.Address.ExactAddress,
			PostalCode:   p.Address.PostalCode,
		}

		response = append(response, parkingDTO)
	}

	return response, nil
}
