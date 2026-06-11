package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	staticdto "github.com/Vanak/Backend/internal/application/dto/static"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/domain/storage"
	"github.com/Vanak/Backend/internal/infrastructure/hashed"
	"gorm.io/gorm"
)

type ParkingService struct {
	unitOfWork ports.UnitOfWork
	cfg        *bootstrap.Config
	logger     logging.Logger
	s3Storage  storage.S3Storage
	hashID     *hashed.HashID
}

func NewParkingService(cfg *bootstrap.Config, unitOfWork ports.UnitOfWork, s3Storage storage.S3Storage, hashedID *hashed.HashID) *ParkingService {
	logger := logging.NewLogger(cfg)

	return &ParkingService{
		logger:     logger,
		cfg:        cfg,
		unitOfWork: unitOfWork,
		s3Storage:  s3Storage,
		hashID:     hashedID,
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
			if fileHeader == nil || fileHeader.Filename == "" || fileHeader.Size == 0 {
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
				continue
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

	err = ps.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		err = rf.ParkingRepository().Delete(parking)

		if err != nil {
			return err
		}

		err := ps.s3Storage.DeleteObject(enums.ParkingDoc, parking.DocumentURL)
		if err != nil {
			return err
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

		return nil

	})

	return err
}

func (ps *ParkingService) UpdateParking(userID uint, parkingInfo parkingdto.UpdateParkingRequest) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	currentHour := uint(now.Hour())
	parkingID := parkingInfo.ParkingID
	return ps.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		parkingRepo := rf.ParkingRepository()
		addressRepo := rf.AddressRepository()

		parking, err := parkingRepo.FindByIDAndUserID(parkingID, userID)
		if err != nil {
			return err
		}

		parking.Title = parkingInfo.Title
		parking.Latitude = parkingInfo.Latitude
		parking.Longitude = parkingInfo.Longitude
		parking.BasePrice = parkingInfo.BasePrice

		existingSlots, err := parkingRepo.GetSlotsByParkingID(parkingID)
		if err != nil {
			return err
		}

		newSlotsGenerated, err := ps.createFirstSlotsForParking(parkingID, parkingInfo.Slots)
		if err != nil {
			return err
		}

		var slotIDsToDelete []uint
		for _, exSlot := range existingSlots {
			found := false
			for _, newSlot := range newSlotsGenerated {
				if exSlot.Date.Format("2006-01-02") == newSlot.Date.Format("2006-01-02") && exSlot.Hour == newSlot.Hour {
					found = true
					break
				}
			}

			if !found {
				isFutureSlot := exSlot.Date.After(today) || (exSlot.Date.Equal(today) && exSlot.Hour >= currentHour)
				if exSlot.IsBooked && isFutureSlot {
					return &exception.ConflictError{
						Item: ps.cfg.Constant.ErrorField.Reserve,
					}
				}
				slotIDsToDelete = append(slotIDsToDelete, exSlot.ID)
			}
		}

		if len(slotIDsToDelete) > 0 {
			if err := parkingRepo.DeleteSlotsWithIDs(slotIDsToDelete); err != nil {
				return err
			}
		}

		var slotsToCreate []entities.AvailabilitySlot
		for _, newSlot := range newSlotsGenerated {
			exists := false
			for _, exSlot := range existingSlots {
				if exSlot.Date.Format("2006-01-02") == newSlot.Date.Format("2006-01-02") && exSlot.Hour == newSlot.Hour {
					exists = true
					break
				}
			}
			if !exists {
				slotsToCreate = append(slotsToCreate, newSlot)
			}
		}

		if len(slotsToCreate) > 0 {
			if err := parkingRepo.CreateSlots(slotsToCreate); err != nil {
				return err
			}
		}

		address, err := addressRepo.GetAddressByParkingID(parkingID)
		if err == nil && address != nil {
			address.CityID = parkingInfo.Address.CityID
			address.PostalCode = parkingInfo.Address.PostalCode
			address.UnitNumber = parkingInfo.Address.UnitNumber
			address.ExactAddress = parkingInfo.Address.ExactAddress

			if err := addressRepo.UpdateAddress(address); err != nil {
				return err
			}
		}

		if parkingInfo.DocumentFile != nil && parkingInfo.DocumentFile.Size > 0 {
			if parking.DocumentURL != "" {
				_ = ps.s3Storage.DeleteObject(enums.ParkingDoc, parking.DocumentURL)
			}

			docKey := ps.cfg.Constant.S3BucketPath.GetParkingDocPath(parking.ID, parkingInfo.DocumentFile.Filename)
			err = ps.s3Storage.UploadFile(enums.ParkingDoc, docKey, parkingInfo.DocumentFile)
			if err != nil {
				return err
			}
			parking.DocumentURL = docKey
		}

		if len(parkingInfo.Images) > 0 {
			for _, exImg := range parking.Images {
				if exImg.ImageURL != "" {
					err := ps.s3Storage.DeleteObject(enums.ParkingPic, exImg.ImageURL)
					if err != nil {
						return err
					}
				}
			}

			var oldImageIDs []uint
			for _, exImg := range parking.Images {
				oldImageIDs = append(oldImageIDs, exImg.ID)
			}
			if len(oldImageIDs) > 0 {
				if err := parkingRepo.DeleteParkingImagesWithIDs(oldImageIDs); err != nil {
					return err
				}
			}

			var newParkingImages []entities.ParkingImage
			for idx, fileHeader := range parkingInfo.Images {
				if fileHeader == nil || fileHeader.Filename == "" || fileHeader.Size == 0 {
					continue
				}

				imageKey := ps.cfg.Constant.S3BucketPath.GetParkingImagePath(parking.ID, idx, fileHeader.Filename)

				err = ps.s3Storage.UploadFile(enums.ParkingPic, imageKey, fileHeader)
				if err != nil {
					return err
				}

				newParkingImages = append(newParkingImages, entities.ParkingImage{
					ImageURL:  imageKey,
					ParkingID: parking.ID,
				})
			}

			if len(newParkingImages) > 0 {
				err = parkingRepo.CreateParkingImages(newParkingImages)
				if err != nil {
					return err
				}
			}
		}

		err = parkingRepo.UpdateParking(parking)
		if err != nil {
			return err
		}

		return nil
	})
}

func (ps *ParkingService) GetParkingDetail(userID uint, parkingID uint) (*parkingdto.ParkingDetailResponse, error) {
	parking, err := ps.unitOfWork.Factory().ParkingRepository().GetDetailByIDAndUserID(parkingID, userID)
	if err != nil {
		return nil, err
	}

	docURL := ""
	if parking.DocumentURL != "" {
		url, err := ps.s3Storage.GetPresignedURL(enums.ParkingDoc, parking.DocumentURL, time.Hour)
		if err == nil {
			docURL = url
		}
	}

	var imageUrls []string
	for _, img := range parking.Images {
		if img.ImageURL != "" {
			url, err := ps.s3Storage.GetPresignedURL(enums.ParkingPic, img.ImageURL, time.Hour)
			if err == nil {
				imageUrls = append(imageUrls, url)
			}
		}
	}

	var attributeIDs []uint
	for _, attr := range parking.Attributes {
		attributeIDs = append(attributeIDs, attr.ID)
	}

	var slotsDTO []parkingdto.SlotDetailResponse
	for _, slot := range parking.Slots {
		slotsDTO = append(slotsDTO, parkingdto.SlotDetailResponse{
			ID:       slot.ID,
			Date:     slot.Date,
			Hour:     slot.Hour,
			IsBooked: slot.IsBooked,
		})
	}

	response := &parkingdto.ParkingDetailResponse{
		ID:          parking.ID,
		Title:       parking.Title,
		Latitude:    parking.Latitude,
		Longitude:   parking.Longitude,
		BasePrice:   uint(parking.BasePrice),
		Status:      parking.Status.String(),
		DocumentURL: docURL,
		ImageURLs:   imageUrls,
		Address: parkingdto.AddressDetailResponse{
			CityID:       parking.Address.CityID,
			PostalCode:   parking.Address.PostalCode,
			UnitNumber:   parking.Address.UnitNumber,
			ExactAddress: parking.Address.ExactAddress,
		},
		Slots:        slotsDTO,
		AttributeIDs: attributeIDs,
	}

	return response, nil
}

func (ps *ParkingService) GetParkingByIDAndUserID(req parkingdto.GetParkingByIDAndUserIDRequest) (*parkingdto.ParkingInfoResponse, error) {
	parkingID, err := ps.hashID.Decode(req.ParkingIDHashed)
	if err != nil {
		return nil, err
	}

	parking, err := ps.unitOfWork.Factory().ParkingRepository().GetParkingByIDAndUserID(req.UserID, parkingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundError(ps.cfg.Constant.ErrorField.Parking)
		}
		return nil, err
	}

	var imageUrls []string
	for _, img := range parking.Images {
		if img.ImageURL != "" {
			url, err := ps.s3Storage.GetPresignedURL(enums.ParkingPic, img.ImageURL, time.Hour)
			if err == nil {
				imageUrls = append(imageUrls, url)
			}
		}
	}

	var attributeIDs []uint
	for _, attr := range parking.Attributes {
		attributeIDs = append(attributeIDs, attr.ID)
	}

	response := &parkingdto.ParkingInfoResponse{
		ID:        req.ParkingIDHashed,
		Title:     parking.Title,
		Longitude: parking.Longitude,
		Latitude:  parking.Latitude,
		BasePrice: parking.BasePrice,
		Status:    parking.Status.String(),
		Address: parkingdto.AddressDetailResponse{
			CityID:       parking.Address.CityID,
			PostalCode:   parking.Address.PostalCode,
			UnitNumber:   parking.Address.UnitNumber,
			ExactAddress: parking.Address.ExactAddress,
		},
		AttributeIDs: attributeIDs,
		ImageURLs:    imageUrls,
	}
	return response, nil
}

func (ps *ParkingService) GetCitiesList() ([]staticdto.StaticItemResponse, error) {
	cities, err := ps.unitOfWork.Factory().AddressRepository().GetAllCities()
	if err != nil {
		return nil, err
	}

	var response []staticdto.StaticItemResponse
	for _, city := range cities {
		response = append(response, staticdto.StaticItemResponse{
			ID:   city.ID,
			Name: city.Name,
		})
	}
	return response, nil
}

func (ps *ParkingService) GetAttributesList() ([]staticdto.StaticItemResponse, error) {
	attributes, err := ps.unitOfWork.Factory().ParkingRepository().GetAllAttributes()
	if err != nil {
		return nil, err
	}

	var response []staticdto.StaticItemResponse
	for _, attr := range attributes {
		response = append(response, staticdto.StaticItemResponse{
			ID:   attr.ID,
			Name: attr.Name,
		})
	}
	return response, nil
}
