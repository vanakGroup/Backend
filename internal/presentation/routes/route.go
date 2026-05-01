package routes

import (
	"github.com/Vanak/Backend/internal/presentation/routes/http/v1"
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	v1 := ginEngine.Group("/v1")
	{
		http.SetupGeneralRoutes(v1)
		http.SetupCustomerRoutes(v1, app)
	}
}
