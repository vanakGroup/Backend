package http

import (
	"github.com/Vanak/Backend/internal/presentation/controller/v1/general"

	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(router *gin.RouterGroup) {

	router.GET("/ping", general.Pong)
}
