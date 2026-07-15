package swagger

import (
	"fmt"

	"github.com/Vanak/Backend/bootstrap"
	docs "github.com/Vanak/Backend/internal/infrastructure/swagger/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwagger(ginEngine *gin.Engine, cfg *bootstrap.Config) {
	docs.SwaggerInfo.Title = "Vanak Backend API"
	docs.SwaggerInfo.Description = "REST API documentation for Vanak core Parking services"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = fmt.Sprintf("vparking.%s", cfg.Constant.Server.TopLevelDomain)
	docs.SwaggerInfo.Schemes = []string{"https"}

	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler((swaggerFiles.Handler)))
}
