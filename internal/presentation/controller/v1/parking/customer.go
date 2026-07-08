package parking

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

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
// @Summary      Create a new parking
// @Description  Create a parking space with details, document, and images
// @Tags         Customer Parking
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        title                formData  string  true   "Parking Title"
// @Param        latitude             formData  number  true   "Latitude (25 to 40)"
// @Param        longitude            formData  number  true   "Longitude (44 to 63)"
// @Param        basePrice            formData  number  true   "Base Price (Greater than 0)"
// @Param        address.cityId       formData  int     true   "City ID"
// @Param        address.postalCode   formData  string  true   "Postal Code"
// @Param        address.unitNumber   formData  int     true   "Unit Number"
// @Param        address.exactAddress formData  string  true   "Exact Address"
// @Param        attributeIds         formData  array   false  "Array of Attribute IDs"
// @Param        slots                formData  string  false  "JSON string of slots array"
// @Param        documentFile         formData  file    true   "Parking legal document file"
// @Param        images               formData  file    false  "Parking gallery images"
// @Success      201      {object}  base.SingleMessageResponse
// @Failure      400      {object}  base.MultipleMessageResponse
// @Router       /parking/crud/create [post]
func (cpc *CustomerParkingController) CreateParking(ctx *gin.Context) {
	type createParkingParams struct {
		Title     string  `form:"title" validate:"required"`
		Latitude  float64 `form:"latitude" validate:"required,gte=25,lte=40"`
		Longitude float64 `form:"longitude" validate:"required,gte=44,lte=63"`
		BasePrice float64 `form:"basePrice" validate:"required,numeric,gt=0"`

		CityID       uint   `form:"address.cityId" validate:"required"`
		PostalCode   string `form:"address.postalCode" validate:"required,ir_postal"`
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

	userID := base.GetID(ctx)
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

	response, err := cpc.parkingService.CreateParking(userID, parkingInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.createParking")
	base.Response(ctx, http.StatusCreated, message, response)
}

// UpdateParking godoc
// @Summary      Update an existing parking
// @Description  Update parking details, location, slots, document and images
// @Tags         Customer Parking
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        parkingId            formData  int     true   "Parking ID to update"
// @Param        title                formData  string  true   "Parking Title"
// @Param        latitude             formData  number  true   "Latitude"
// @Param        longitude            formData  number  true   "Longitude"
// @Param        basePrice            formData  number  true   "Base Price"
// @Param        address.cityId       formData  int     true   "City ID"
// @Param        address.postalCode   formData  string  true   "Postal Code"
// @Param        address.unitNumber   formData  int     true   "Unit Number"
// @Param        address.exactAddress formData  string  true   "Exact Address"
// @Param        attributeIds         formData  array   false  "Array of Attribute IDs"
// @Param        slots                formData  string  false  "JSON string of slots array"
// @Param        documentFile         formData  file    true   "New document file"
// @Param        images               formData  file    false  "New complete package of gallery images"
// @Success      201      {object}  base.SingleMessageResponse
// @Failure      400      {object}  base.MultipleMessageResponse
// @Router       /parking/crud/update [put]
func (cpc *CustomerParkingController) UpdateParking(ctx *gin.Context) {
	type updateParkingParams struct {
		ParkingID string  `form:"parkingId" validate:"required,min=1"`
		Title     string  `form:"title" validate:"required"`
		Latitude  float64 `form:"latitude" validate:"required,gte=25,lte=40"`
		Longitude float64 `form:"longitude" validate:"required,gte=44,lte=63"`
		BasePrice float64 `form:"basePrice" validate:"required,numeric,gt=0"`

		CityID       uint   `form:"address.cityId" validate:"required"`
		PostalCode   string `form:"address.postalCode" validate:"required,ir_postal"`
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

	userID := base.GetID(ctx)

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

	err = cpc.parkingService.UpdateParking(userID, updateParkingInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.updateParking")
	base.Response(ctx, http.StatusCreated, message, nil)
}

// DeleteParking godoc
// @Summary      Delete a parking
// @Description  Delete a parking by its ID checking ownership
// @Tags         Customer Parking
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Parking ID"
// @Success      200  {object}  base.SingleMessageResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /parking/crud/delete/{id} [delete]
func (cpc *CustomerParkingController) DeleteParking(ctx *gin.Context) {
	type DeleteParkingParam struct {
		IDHashed string `uri:"id" validate:"required,min=1"`
	}

	params := base.Validated[DeleteParkingParam](ctx)

	userID := base.GetID(ctx)

	err := cpc.parkingService.DeleteParking(userID, params.IDHashed)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteParking")
	base.Response(ctx, http.StatusOK, message, nil)
}

// GetParkingDetail godoc
// @Summary      Get parking details for edit form
// @Description  Retrieve complete parking info with relations and presigned URLs by ID
// @Tags         Customer Parking
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Parking ID"
// @Success      200  {object}  parkingdto.ParkingDetailResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /parking/crud/parkingDetail/{id} [get]
func (cpc *CustomerParkingController) GetParkingDetail(ctx *gin.Context) {
	type GetParkingDetailParam struct {
		IDHashed string `uri:"id" validate:"required,min=1"`
	}

	params := base.Validated[GetParkingDetailParam](ctx)

	userID := base.GetID(ctx)

	parking, err := cpc.parkingService.GetParkingDetail(userID, params.IDHashed)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.detailParking")
	base.Response(ctx, http.StatusOK, message, parking)
}

// SearchParking godoc
// @Summary      Search parkings
// @Description  Search available parkings with optional filters. Pagination: send page + pageSize, offset is calculated here.
// @Tags         Parking Search
// @Accept       json
// @Produce      json
// @Param        request  body  object  false  "Search filters"
// @Success      200  {object}  parkingdto.ParkingSearchResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /parking/search [post]
func (cpc *CustomerParkingController) SearchParking(ctx *gin.Context) {
	type searchParkingParams struct {
		Page     int `json:"page"     validate:"omitempty,min=1"`
		PageSize int `json:"pageSize" validate:"omitempty,min=1,max=100"`

		SortBy  string `json:"sortBy"  validate:"omitempty,oneof=price"`
		SortAsc bool   `json:"sortAsc"`

		MinPrice *float64 `json:"minPrice"  validate:"omitempty,gt=0"`
		MaxPrice *float64 `json:"maxPrice"  validate:"omitempty,gt=0"`

		CityID  *uint   `json:"cityId"  validate:"omitempty"`
		Address *string `json:"address" validate:"omitempty,min=2"`

		Latitude  *float64 `json:"latitude"  validate:"omitempty"`
		Longitude *float64 `json:"longitude" validate:"omitempty"`
		Radius    *float64 `json:"radius"    validate:"omitempty,gt=0"`

		DateRaw   string `json:"date"      validate:"omitempty"`
		StartHour uint   `json:"startHour" validate:"omitempty,max=23"`
		EndHour   uint   `json:"endHour"   validate:"omitempty,max=24"`

		AttributeIDs []uint `json:"attributeIds" validate:"omitempty"`
	}

	params := base.Validated[searchParkingParams](ctx)

	offset, limit := base.GetOffsetLimit(
		params.Page,
		params.PageSize,
		cpc.cfg.Constant.Pagination.DefaultPage,
		cpc.cfg.Constant.Pagination.DefaultPageSize,
	)

	sortColumn := ""
	if params.SortBy == "price" {
		sortColumn = "base_price"
	}

	req := parkingdto.ParkingSearchRequest{
		Limit:        limit,
		Offset:       offset,
		SortBy:       sortColumn,
		SortAsc:      params.SortAsc,
		MinPrice:     params.MinPrice,
		MaxPrice:     params.MaxPrice,
		CityID:       params.CityID,
		Address:      params.Address,
		Latitude:     params.Latitude,
		Longitude:    params.Longitude,
		Radius:       params.Radius,
		StartHour:    params.StartHour,
		EndHour:      params.EndHour,
		AttributeIDs: params.AttributeIDs,
	}

	if params.DateRaw != "" {
		parsed, err := time.Parse("2006-01-02", params.DateRaw)
		if err != nil {
			panic(exception.NewBinding(err))
		}
		req.Date = &parsed
	}

	data, total, err := cpc.parkingService.SearchParking(req)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.searchParking")
	base.Response(ctx, http.StatusOK, message, base.NewPaginatedResponse(data, total, offset, limit))
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

// GetCities godoc
// @Summary      Get all cities
// @Description  Retrieve a list of all cities with their IDs for dropdowns
// @Tags         Static Data
// @Produce      json
// @Success      200  {array}   staticdto.StaticItemResponse
// @Failure      500  {object}  base.MultipleMessageResponse
// @Router       /parking/static/cities [get]
func (cpc *CustomerParkingController) GetCities(ctx *gin.Context) {
	cities, err := cpc.parkingService.GetCitiesList()
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getCities")
	base.Response(ctx, http.StatusOK, message, cities)
}

// GetAttributes godoc
// @Summary      Get all parking attributes
// @Description  Retrieve a list of all available parking attributes/features with their IDs
// @Tags         Static Data
// @Produce      json
// @Success      200  {array}   staticdto.StaticItemResponse
// @Failure      500  {object}  base.MultipleMessageResponse
// @Router       /parking/static/attributes [get]
func (cpc *CustomerParkingController) GetAttributes(ctx *gin.Context) {
	attributes, err := cpc.parkingService.GetAttributesList()
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getAttributes")
	base.Response(ctx, http.StatusOK, message, attributes)
}

// GetParkingInfo godoc
// @Summary      Get parking Info
// @Description  Retrieve complete parking info
// @Tags         Customer Parking
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Parking ID"
// @Success      200  {object}  parkingdto.ParkingDetailResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /parking/parkingInfo/{id} [get]
func (cpc *CustomerParkingController) GetParkingInfo(ctx *gin.Context) {
	type GetParkingInfoParam struct {
		IDHashed string `uri:"id" validate:"required,min=1"`
	}

	params := base.Validated[GetParkingInfoParam](ctx)
	userID := base.GetID(ctx)

	req := parkingdto.GetParkingByIDAndUserIDRequest{
		UserID:          userID,
		ParkingIDHashed: params.IDHashed,
	}

	res, err := cpc.parkingService.GetParkingByIDAndUserID(req)

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getParkingInfo")
	base.Response(ctx, http.StatusOK, message, res)
}

// GetUserParkings godoc
// @Summary      Get User parkings
// @Description  Retrieve complete User parkings
// @Tags         Customer Parking
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  parkingdto.ParkingDetailResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /parking [get]
func (cpc *CustomerParkingController) GetUserParkings(ctx *gin.Context) {
	userID := base.GetID(ctx)

	res, err := cpc.parkingService.GetUserParkings(userID)

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getUserParkings")
	base.Response(ctx, http.StatusOK, message, res)
}
