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
	}
}
