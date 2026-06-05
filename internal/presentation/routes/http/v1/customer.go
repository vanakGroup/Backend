package http

import (
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(router *gin.RouterGroup, app *wire.Application) {
	parkingCrudGroup := router.Group("/parking/crud")
	parkingCrudGroup.Use(app.Middlewares.Authentication.AuthMiddleware())
	{
		parkingCrudGroup.POST("/create", app.CustomerParkingController.CreateParking)
		parkingCrudGroup.DELETE("/delete/:id", app.CustomerParkingController.DeleteParking)
		parkingCrudGroup.GET("/parkingDetail/:id", app.CustomerParkingController.GetParkingDetail)
		parkingCrudGroup.PUT("/update", app.CustomerParkingController.UpdateParking)
	}
}
