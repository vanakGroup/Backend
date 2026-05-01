package http

import (
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(router *gin.RouterGroup, app *wire.Application) {
	AddGroup := router.Group("/math")
	{
		AddGroup.POST("add/:first/:second", app.MathController.Adder)
	}
}
