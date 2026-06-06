package vehicle

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/Vanak/Backend/bootstrap"
	vehicledto "github.com/Vanak/Backend/internal/application/dto/vehicle"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type CustomerVehicleController struct {
	cfg            *bootstrap.Config
	vehicleService usecase.VehicleService
}

func NewVehicleCustomerController(cfg *bootstrap.Config, vehicleService usecase.VehicleService) *CustomerVehicleController {
	return &CustomerVehicleController{
		cfg:            cfg,
		vehicleService: vehicleService,
	}
}

// CreateVehicle godoc
// @Summary      Register a new customer vehicle
// @Description  Receive text information (license plate and car model) along with the uploaded vehicle registration card image.
// @Tags         Customer Vehicle
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        license_plate      formData  string  true  "Vehicle license plate (e.g., 12B345-67)"
// @Param        car_model          formData  string  true  "Car model (e.g., Peugeot Pars)"
// @Param        registration_card  formData  file    true  "Physical vehicle registration card image (JPG/PNG)"
// @Success      200  {object}  base.SingleMessageResponse "Successful registration message"
// @Failure      400  {object}  base.SingleMessageResponse "Validation error (missing required fields or file)"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error (file upload or database issue)"
// @Router       /vehicle [post]
func (cvc CustomerVehicleController) CreateVehicle(ctx *gin.Context) {

	type VehiclesFilledParams struct {
		LicensePlate     string                `form:"license_plate" validate:"required"`
		CarModel         string                `form:"car_model" validate:"required"`
		RegistrationCard *multipart.FileHeader `form:"registration_card" validate:"required"`
	}
	params := base.Validated[VehiclesFilledParams](ctx)

	userIDAny, exists := ctx.Get(cvc.cfg.Constant.Context.UserID)
	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cvc.cfg.Constant.ErrorField.User, cvc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userId := userIDAny.(uint)

	vehiclesInfo := vehicledto.VehicleDTORequest{
		LicensePlate: params.LicensePlate,
		CarModel:     params.CarModel,
		CardFile:     params.RegistrationCard,
		CardFileName: params.RegistrationCard.Filename,
		UserId:       userId,
	}

	err := cvc.vehicleService.RegisterVehicle(ctx.Request.Context(), vehiclesInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cvc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.registerVehicle")
	base.Response(ctx, 200, message, nil)

}

// GetAllVehicles godoc
// @Summary      Get list of customer vehicles
// @Description  Retrieve a summary list of all vehicles belonging to the authenticated customer.
// @Tags         Customer Vehicle
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   vehicledto.VehicleListDTOResponse "Vehicles list fetched successfully"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /vehicle [get]
func (cvc *CustomerVehicleController) GetAllVehicles(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(cvc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cvc.cfg.Constant.ErrorField.User, cvc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userId := userIDAny.(uint)
	vehiclesList, err := cvc.vehicleService.GetAllVehicles(userId)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cvc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getVehiclesList")
	base.Response(ctx, 200, message, vehiclesList)
}

// GetVehicle godoc
// @Summary      Get customer vehicle details
// @Description  Retrieve the details of a specific vehicle, including the presigned URL for the registration card.
// @Tags         Customer Vehicle
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Vehicle ID"
// @Success      200  {object}  vehicledto.VehicleDTOResponse "Vehicle details successfully retrieved"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid ID format"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      403  {object}  base.SingleMessageResponse "Forbidden (Vehicle does not belong to the user)"
// @Failure      404  {object}  base.SingleMessageResponse "Vehicle not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /vehicle/{id} [get]
func (cvc *CustomerVehicleController) GetVehicle(ctx *gin.Context) {
	vehicleIDString := ctx.Param("id")
	userIDAny, exists := ctx.Get(cvc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cvc.cfg.Constant.ErrorField.User, cvc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	vehicleID64, err := strconv.ParseUint(vehicleIDString, 10, 0)
	if err != nil {
		formatError := exception.NewFormatError()
		formatError.Add(cvc.cfg.Constant.VehicleFields.ID, cvc.cfg.Constant.ErrorTag.WrongFormat, cvc.cfg.Constant.ErrorField.VehicleId)
		panic(formatError)
	}

	userId := userIDAny.(uint)
	vehicleID := uint(vehicleID64)
	vehicleInfo, informationErr := cvc.vehicleService.GetVehicle(vehicleID, userId)

	if informationErr != nil {
		panic(informationErr)
	}

	trans := base.GetTranslator(ctx, cvc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getVehicle")
	base.Response(ctx, 200, message, vehicleInfo)

}

// DeleteVehicle godoc
// @Summary      Delete customer vehicle
// @Description  Remove a vehicle belonging to the authenticated customer. Currently, this soft-deletes the database record and permanently removes the image from S3.
// @Tags         Customer Vehicle
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Vehicle ID"
// @Success      200  {object}  base.SingleMessageResponse "Vehicle successfully deleted"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid ID format"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      404  {object}  base.SingleMessageResponse "Vehicle not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /vehicle/{id} [delete]
func (cvc *CustomerVehicleController) DeleteVehicle(ctx *gin.Context) {
	vehicleIDString := ctx.Param("id")
	userIDAny, exists := ctx.Get(cvc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cvc.cfg.Constant.ErrorField.User, cvc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	vehicleID64, err := strconv.ParseUint(vehicleIDString, 10, 0)
	if err != nil {
		formatError := exception.NewFormatError()
		formatError.Add(cvc.cfg.Constant.VehicleFields.ID, cvc.cfg.Constant.ErrorTag.WrongFormat, cvc.cfg.Constant.ErrorField.VehicleId)
		panic(formatError)
	}

	userId := userIDAny.(uint)
	vehicleID := uint(vehicleID64)

	err = cvc.vehicleService.DeleteVehicle(vehicleID, userId)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cvc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteVehicle")
	base.Response(ctx, 200, message, nil)

}
