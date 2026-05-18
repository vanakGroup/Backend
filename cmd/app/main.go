package main

import (
	"log"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/infrastructure/validation"
	"github.com/Vanak/Backend/internal/presentation/middleware"
	"github.com/Vanak/Backend/internal/presentation/routes"
	"github.com/Vanak/Backend/wire"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	cfg := bootstrap.Run()

	app, err := wire.InitializeApplication(cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	gin.DisableConsoleColor()
	ginEngine := gin.Default()

	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterValidation("ir_phone", validation.IranianMobileValidator, true)
	}

	ginEngine.Use(middleware.TranslatorMiddleware(app.Translator))
	ginEngine.Use(app.Middlewares.Recovery.Recovery)

	if app.Database.DB == nil {
		log.Fatal("Database instance is nil")
	}

	routes.Run(ginEngine, app)

	log.Println("🚀 Server is running on port 8080...")
	ginEngine.Run()
}
