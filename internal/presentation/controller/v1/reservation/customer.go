package reservation

import (
	"fmt"

	"github.com/Vanak/Backend/bootstrap"
	reservationdto "github.com/Vanak/Backend/internal/application/dto/reservation"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type CustomerReservationController struct {
	cfg                *bootstrap.Config
	reservationService usecase.ReservationService
}

func NewCustomerReservationController(
	cfg *bootstrap.Config,
	reservationService usecase.ReservationService,
) *CustomerReservationController {
	return &CustomerReservationController{
		cfg:                cfg,
		reservationService: reservationService,
	}
}

// CreateReservation godoc
// @Summary      Create a new parking reservation
// @Description  Register a new parking reservation by providing hashed parking and vehicle IDs, along with start and end times as Unix timestamps.
// @Tags         Customer Reservation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object{parking_id=string,vehicle_id=string,start_time=integer,end_time=integer} true "Reservation payload (Timestamps must be > 0 and end_time > start_time)" example(object{parking_id="hash123",vehicle_id="hash456",start_time=1718000000,end_time=1718010800})
// @Success      200  {object}  base.SingleMessageResponse "Successful reservation message"
// @Failure      400  {object}  base.SingleMessageResponse "Validation error (missing fields, invalid time range, or invalid IDs)"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      409  {object}  base.SingleMessageResponse "Conflict error (Slots are unavailable or already booked)"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /reservation [post]
func (crc *CustomerReservationController) CreateReservation(ctx *gin.Context) {
	type ReservationFieldsParam struct {
		ParkingID string `json:"parking_id" binding:"required"`
		VehicleID string `json:"vehicle_id" binding:"required"`
		StartTime int64  `json:"start_time" binding:"required,gt=0"`
		EndTime   int64  `json:"end_time" binding:"required,gtfield=StartTime"`
	}

	params := base.Validated[ReservationFieldsParam](ctx)

	userIDAny, exists := ctx.Get(crc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", crc.cfg.Constant.ErrorField.User, crc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)

	reservationInfo := reservationdto.ReservationDTORequest{
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
		UserID:    userID,
		ParkingID: params.ParkingID,
		VehicleID: params.VehicleID,
	}

	err := crc.reservationService.CreateReservation(reservationInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, crc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.registerReservation")
	base.Response(ctx, 200, message, nil)
}

// GetParkingReservationList godoc
// @Summary      Get user's reservations for a specific parking
// @Description  Retrieve a list of all reservations made by the authenticated customer for a specific parking facility using its hashed ID.
// @Tags         Parking Reservation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        parking_id   path      string  true  "Hashed Parking ID"
// @Success      200  {object}  base.SingleMessageResponse "List of reservations successfully retrieved"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid parking ID format"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      404  {object}  base.SingleMessageResponse "Parking not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /parking/{parking_id}/reservation/ [get]
func (crc *CustomerReservationController) GetParkingReservationList(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(crc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", crc.cfg.Constant.ErrorField.User, crc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)

	parkingID := ctx.Param("parking_id")

	reservationList, err := crc.reservationService.GetReservationListForParking(parkingID, userID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, crc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getReservationList")
	base.Response(ctx, 200, message, reservationList)
}

// GetParkingReservationDetails godoc
// @Summary      Get details of a specific parking reservation
// @Description  Retrieve detailed information of a specific reservation inside a specific parking lot using their hashed IDs.
// @Tags         Parking Reservation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        parking_id   path      string  true  "Hashed Parking ID"
// @Param        reserve_id   path      string  true  "Hashed Reservation ID"
// @Success      200  {object}  base.SingleMessageResponse "Reservation details successfully retrieved"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid IDs provided"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      404  {object}  base.SingleMessageResponse "Reservation or Parking not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /parking/{parking_id}/reservation/{reserve_id} [get]
func (crc *CustomerReservationController) GetParkingReservationDetails(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(crc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", crc.cfg.Constant.ErrorField.User, crc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	parkingID := ctx.Param("parking_id")
	reserveID := ctx.Param("reserve_id")
	userID := userIDAny.(uint)

	reservation, err := crc.reservationService.GetReservationDetailForParking(reserveID, parkingID, userID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, crc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getReservationInfo")
	base.Response(ctx, 200, message, reservation)
}

// GetVehicleReservationList godoc
// @Summary      Get user's reservations for a specific vehicle
// @Description  Retrieve a list of all reservations made by the authenticated customer for a specific vehicle using its hashed ID.
// @Tags         Vehicle Reservation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Hashed Vehicle ID"
// @Success      200  {object}  base.SingleMessageResponse "List of vehicle reservations successfully retrieved"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid vehicle ID format or missing ID"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      404  {object}  base.SingleMessageResponse "Vehicle not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /vehicle/{id}/reservation [get]
func (crc *CustomerReservationController) GetVehicleReservationList(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(crc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", crc.cfg.Constant.ErrorField.User, crc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)

	vehicleID := ctx.Param("id")

	reservationList, err := crc.reservationService.GetReservationListForVehicle(vehicleID, userID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, crc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getReservationList")
	base.Response(ctx, 200, message, reservationList)
}

// GetVehicleReservationDetails godoc
// @Summary      Get details of a specific reservation for a vehicle
// @Description  Retrieve detailed information of a specific reservation for a given vehicle using their hashed IDs.
// @Tags         Vehicle Reservation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Hashed Vehicle ID"
// @Param        reserve_id   path      string  true  "Hashed Reservation ID"
// @Success      200  {object}  base.SingleMessageResponse "Reservation details successfully retrieved"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid IDs provided"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      404  {object}  base.SingleMessageResponse "Reservation or Vehicle not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /vehicle/{id}/reservation/{reserve_id} [get]
func (crc *CustomerReservationController) GetVehicleReservationDetails(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(crc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", crc.cfg.Constant.ErrorField.User, crc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)

	vehicleID := ctx.Param("id")
	reserveID := ctx.Param("reserve_id")

	reservation, err := crc.reservationService.GetReservationDetailForVehicle(reserveID, vehicleID, userID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, crc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getReservationInfo")
	base.Response(ctx, 200, message, reservation)
}

// VerifySettlement godoc
// @Summary      Verify and settle a reservation
// @Description  Settle the payment for a specific reservation and transfer funds to the parking owner.
// @Tags         Reservation Settlement
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        reserve_id   path      string  true  "Hashed Reservation ID"
// @Param        parking_id   path      string  true  "Hashed Parking ID"
// @Success      200  {object}  base.SingleMessageResponse "Settlement successful"
// @Failure      400  {object}  base.SingleMessageResponse "Invalid IDs or Insufficient funds"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /parking/{parking_id}/reservation/{reserve_id}/settle [post]
func (crc *CustomerReservationController) Settle(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(crc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", crc.cfg.Constant.ErrorField.User, crc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)
	reserveID := ctx.Param("reserve_id")
	parkingID := ctx.Param("parking_id")

	err := crc.reservationService.VerifySettlement(reserveID, parkingID, userID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, crc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.verifySettlement")
	base.Response(ctx, 200, message, nil)
}
