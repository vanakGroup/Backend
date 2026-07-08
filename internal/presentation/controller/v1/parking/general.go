package parking

import (
	"net/http"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	parkingdto "github.com/Vanak/Backend/internal/application/dto/parking"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type GeneralParkingController struct {
	parkingService usecase.ParkingService
	cfg            *bootstrap.Config
}

func NewGeneralParkingController(
	cfg *bootstrap.Config,
	parkingService usecase.ParkingService,
) *GeneralParkingController {
	return &GeneralParkingController{
		cfg:            cfg,
		parkingService: parkingService,
	}
}

// GetAllParkings godoc
// @Summary      Get all parkings
// @Description  Retrieve paginated list of all parkings without authentication
// @Tags         Public Parking
// @Produce      json
// @Param        page      query  int  false  "Page number (default 1)"
// @Param        pageSize  query  int  false  "Items per page (default 10, max 100)"
// @Success      200  {object}  parkingdto.ParkingSearchResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /Generalparking [get]
func (gpc *GeneralParkingController) GetAllParkings(ctx *gin.Context) {
	type params struct {
		Page     int `json:"page"     validate:"omitempty,min=1"`
		PageSize int `json:"pageSize" validate:"omitempty,min=1,max=100"`
	}

	p := base.Validated[params](ctx)

	offset, limit := base.GetOffsetLimit(
		p.Page,
		p.PageSize,
		gpc.cfg.Constant.Pagination.DefaultPage,
		gpc.cfg.Constant.Pagination.DefaultPageSize,
	)

	data, total, err := gpc.parkingService.GetAllParkings(limit, offset)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, gpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getUserParkings")
	base.Response(ctx, http.StatusOK, message, base.NewPaginatedResponse(data, total, offset, limit))
}

// GetParkingPublicInfo godoc
// @Summary      Get public parking info
// @Description  Retrieve general parking information by ID without authentication
// @Tags         Public Parking
// @Produce      json
// @Param        id  path  string  true  "Parking ID"
// @Success      200  {object}  parkingdto.ParkingInfoResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /Generalparking/info/{id} [get]
func (gpc *GeneralParkingController) GetParkingPublicInfo(ctx *gin.Context) {
	type pathParams struct {
		IDHashed string `uri:"id" validate:"required,min=1"`
	}

	p := base.Validated[pathParams](ctx)

	res, err := gpc.parkingService.GetParkingPublicInfo(p.IDHashed)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, gpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getParkingInfo")
	base.Response(ctx, http.StatusOK, message, res)
}

// SearchParking godoc
// @Summary      Search parkings
// @Description  Search available parkings with optional filters. Pagination: send page + pageSize, offset is calculated here.
// @Tags         Public Parking
// @Accept       json
// @Produce      json
// @Param        request  body  object  false  "Search filters"
// @Success      200  {object}  parkingdto.ParkingSearchResponse
// @Failure      400  {object}  base.MultipleMessageResponse
// @Router       /Generalparking/search [post]
func (gpc *GeneralParkingController) SearchParking(ctx *gin.Context) {
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

	var params searchParkingParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		panic(err)
	}

	offset, limit := base.GetOffsetLimit(
		params.Page,
		params.PageSize,
		gpc.cfg.Constant.Pagination.DefaultPage,
		gpc.cfg.Constant.Pagination.DefaultPageSize,
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

	data, total, err := gpc.parkingService.SearchParking(req)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, gpc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.searchParking")
	base.Response(ctx, http.StatusOK, message, base.NewPaginatedResponse(data, total, offset, limit))
}
