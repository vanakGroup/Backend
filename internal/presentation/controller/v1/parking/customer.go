package parking

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Vanak/Backend/bootstrap"
	addressdto "github.com/Vanak/Backend/internal/application/dto/address"
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type CustomerParkingController struct {
	parkingService usecase.ParkingService
	cfg            *bootstrap.Config
}

func NewCustomerParkingController(
	cfg *bootstrap.Config,
	parkingService usecase.ParkingService,
) *CustomerParkingController {
	return &CustomerParkingController{
		cfg:            cfg,
		parkingService: parkingService,
	}
}

// CreateParking godoc
// @Summary      Creating Parking
// @Description  Create the parking for user
// @Tags         Authentication
// @Param        request  body      object{phone=string}  true  "User's mobile phone number"  example(object{phone="09123456789"})
// @Success      200      {object}  base.SingleMessageResponse
// @Failure      400      {object}  base.MultipleMessageResponse
// @Router       /parking/crud/create [post]
func (cpc *CustomerParkingController) CreateParking(ctx *gin.Context) {
	type createParkingParams struct {
		Title     string  `form:"title" validate:"required"`
		Latitude  float64 `form:"latitude" validate:"required,gte=25,lte=40"`
		Longitude float64 `form:"longitude" validate:"required,gte=44,lte=63"`
		BasePrice float64 `form:"basePrice" validate:"required,numeric,gt=0"`

		CityID       uint   `form:"address.cityId" validate:"required"`
		PostalCode   string `form:"address.postalCode" validate:"required"`
		UnitNumber   uint   `form:"address.unitNumber" validate:"required"`
		ExactAddress string `form:"address.exactAddress" validate:"required"`

		AttributeIDs []uint `form:"attributeIds" validate:"omitempty"`
		SlotsRaw     string `form:"slots" validate:"omitempty"`

		DocumentFile *multipart.FileHeader   `form:"documentFile" validate:"required"`
		Images       []*multipart.FileHeader `form:"images" validate:"omitempty"`
	}

	params := base.Validated[createParkingParams](ctx)

	slots, err := parseSlotsRaw(params.SlotsRaw)
	if err != nil {
		NewBindingError := exception.NewBinding(err)
		panic(NewBindingError)
	}

	userID, _ := ctx.Get(cpc.cfg.Constant.Context.UserID)
	fmt.Println("the User id : ", userID)
	parkingInfo := parkingdto.CreateParkingRequest{
		Title:     params.Title,
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
		BasePrice: params.BasePrice,
		Address: addressdto.AddressRequest{
			CityID:       params.CityID,
			PostalCode:   params.PostalCode,
			UnitNumber:   params.UnitNumber,
			ExactAddress: params.ExactAddress,
		},
		AttributeIDs: params.AttributeIDs,
		Slots:        slots,
		DocumentFile: params.DocumentFile,
		Images:       params.Images,
	}

	response, err := cpc.parkingService.CreateParking(userID.(uint), parkingInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.createParking")
	base.Response(ctx, http.StatusCreated, message, response)
}

func (cpc *CustomerParkingController) UpdateParking(ctx *gin.Context) {
	type updateParkingParams struct {
		ParkingID uint    `form:"parkingId" validate:"required"`
		Title     string  `form:"title" validate:"required"`
		Latitude  float64 `form:"latitude" validate:"required,gte=25,lte=40"`
		Longitude float64 `form:"longitude" validate:"required,gte=44,lte=63"`
		BasePrice float64 `form:"basePrice" validate:"required,numeric,gt=0"`

		CityID       uint   `form:"address.cityId" validate:"required"`
		PostalCode   string `form:"address.postalCode" validate:"required"`
		UnitNumber   uint   `form:"address.unitNumber" validate:"required"`
		ExactAddress string `form:"address.exactAddress" validate:"required"`

		AttributeIDs []uint `form:"attributeIds" validate:"omitempty"`
		SlotsRaw     string `form:"slots" validate:"omitempty"`

		DocumentFile *multipart.FileHeader   `form:"documentFile" validate:"required"`
		Images       []*multipart.FileHeader `form:"images" validate:"omitempty"`
	}

	params := base.Validated[updateParkingParams](ctx)

	slots, err := parseSlotsRaw(params.SlotsRaw)
	if err != nil {
		NewBindingError := exception.NewBinding(err)
		panic(NewBindingError)

	}

	userID, _ := ctx.Get(cpc.cfg.Constant.Context.UserID)

	updateParkingInfo := parkingdto.UpdateParkingRequest{
		ParkingID: params.ParkingID,
		Title:     params.Title,
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
		BasePrice: params.BasePrice,
		Address: addressdto.UpdateAddressRequest{
			CityID:       params.CityID,
			PostalCode:   params.PostalCode,
			UnitNumber:   params.UnitNumber,
			ExactAddress: params.ExactAddress,
		},
		AttributeIDs: params.AttributeIDs,
		Slots:        slots,
		DocumentFile: params.DocumentFile,
		Images:       params.Images,
	}

	err = cpc.parkingService.UpdateParking(userID.(uint), updateParkingInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.updateParking")
	base.Response(ctx, http.StatusCreated, message, nil)
}

func (cpc *CustomerParkingController) DeleteParking(ctx *gin.Context) {
	parkingIDStr := ctx.Param("id")
	var parkingID uint
	_, err := fmt.Sscanf(parkingIDStr, "%d", &parkingID)
	if err != nil || parkingID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "شناسه پارکینگ معتبر نیست"})
		return
	}

	userContextValue, exists := ctx.Get(cpc.cfg.Constant.Context.UserID)
	if !exists || userContextValue == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "کاربر احراز هویت نشده است"})
		return
	}
	userID := userContextValue.(uint)

	err = cpc.parkingService.DeleteParking(userID, parkingID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteParking")
	base.Response(ctx, http.StatusOK, message, nil)
}

// GetUserParkings godoc
// @Summary      Get Owner Parkings
// @Description  Retrieve all parkings belonging to the authenticated user with presigned asset URLs
// @Tags         Authentication
// @Success      200  {array}   parkingdto.ParkingResponse
// @Failure      401  {object}  base.MultipleMessageResponse
// @Router       /v1/parking/crud/list [get]
func (cpc *CustomerParkingController) GetParkingDetail(ctx *gin.Context) {
	parkingIDStr := ctx.Param("id")
	parkingID, _ := strconv.Atoi(parkingIDStr)
	userID, _ := ctx.Get(cpc.cfg.Constant.Context.UserID)

	parking, err := cpc.parkingService.GetParkingDetail(userID.(uint), uint(parkingID))
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.detailParking")
	base.Response(ctx, http.StatusOK, message, parking)
}

func parseSlotsRaw(slotsRaw string) ([]parkingdto.SlotRequest, error) {
	var slots []parkingdto.SlotRequest
	if slotsRaw == "" {
		return slots, nil
	}

	if err := json.Unmarshal([]byte(slotsRaw), &slots); err != nil {
		return nil, err
	}

	return slots, nil
}
