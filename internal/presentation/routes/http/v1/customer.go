package http

import (
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(router *gin.RouterGroup, app *wire.Application) {
	parkingGroup := router.Group("/parking")
	parkingGroup.Use(app.Middlewares.Authentication.AuthMiddleware()).Use(app.Middlewares.Authentication.UserExistMiddleware())
	{
		parkingGroup.GET("", app.CustomerParkingController.GetUserParkings)
		parkingGroup.GET("/info/:id", app.CustomerParkingController.GetParkingInfo)
	}

	parkingCrudGroup := parkingGroup.Group("/crud")
	parkingCrudGroup.Use(app.Middlewares.Authentication.AuthMiddleware()).Use(app.Middlewares.Authentication.UserExistMiddleware())
	{
		parkingCrudGroup.POST("/create", app.CustomerParkingController.CreateParking)
		parkingCrudGroup.DELETE("/delete/:id", app.CustomerParkingController.DeleteParking)
		parkingCrudGroup.GET("/parkingDetail/:id", app.CustomerParkingController.GetParkingDetail)
		parkingCrudGroup.PUT("/update", app.CustomerParkingController.UpdateParking)
	}

	staticGroup := router.Group("/parking/static")
	{
		staticGroup.GET("/cities", app.CustomerParkingController.GetCities)
		staticGroup.GET("/attributes", app.CustomerParkingController.GetAttributes)
	}
	profileGroup := router.Group("/profile")

	profileGroup.Use(app.Middlewares.Authentication.AuthMiddleware())
	{
		profileGroup.PUT("", app.CustomerUserController.UpdateProfile)
		profileGroup.GET("", app.CustomerUserController.GetProfile)
	}

	vehicleGroup := router.Group("/vehicle")
	vehicleGroup.Use(app.Middlewares.Authentication.AuthMiddleware())
	{
		vehicleGroup.POST("", app.CustomerVehicleController.CreateVehicle)
		vehicleGroup.GET("", app.CustomerVehicleController.GetAllVehicles)
		vehicleGroup.GET("/:id", app.CustomerVehicleController.GetVehicle)
		vehicleGroup.DELETE("/:id", app.CustomerVehicleController.DeleteVehicle)
	}

	walletGroup := router.Group("/wallet")
	walletGroup.Use(app.Middlewares.Authentication.AuthMiddleware())
	{
		walletGroup.GET("", app.CustomerWalletController.GetWallet)
		walletGroup.POST("/deposit", app.CustomerWalletController.Deposit)
		walletGroup.GET("/transactions", app.CustomerWalletController.GetLastTransactions)
	}
}
