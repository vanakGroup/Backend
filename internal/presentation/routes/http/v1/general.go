package http

import (
	"github.com/Vanak/Backend/internal/presentation/controller/v1/general"
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(router *gin.RouterGroup, app *wire.Application) {

	router.GET("/ping", general.Pong)
	router.GET("/send-otp", app.GeneralUserController.SendOtp)
}
