package http

import (
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(router *gin.RouterGroup, app *wire.Application) {
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
}
