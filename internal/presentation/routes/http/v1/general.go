package http

import (
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(router *gin.RouterGroup, app *wire.Application) {

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/send-otp", app.GeneralUserController.SendOtp)
		authGroup.POST("/verify-otp", app.GeneralUserController.VerifyOTP)
		authGroup.POST("/login", app.GeneralUserController.Login)
	}

	parkingGroup := router.Group("/Generalparking")
	{
		parkingGroup.GET("", app.GeneralParkingController.GetAllParkings)
		parkingGroup.GET("/info/:id", app.GeneralParkingController.GetParkingPublicInfo)
		parkingGroup.POST("/search", app.GeneralParkingController.SearchParking)
	}

}
