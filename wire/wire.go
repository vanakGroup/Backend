//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/application/service"
	mathService "github.com/Vanak/Backend/internal/application/service/math"
	"github.com/Vanak/Backend/internal/application/usecase"
	domainLocalization "github.com/Vanak/Backend/internal/domain/localization"
	redisDomain "github.com/Vanak/Backend/internal/domain/repository/redis"
	"github.com/Vanak/Backend/internal/infrastructure/database"
	"github.com/Vanak/Backend/internal/infrastructure/localization"
	redisRepo "github.com/Vanak/Backend/internal/infrastructure/repository/redis"
	"github.com/Vanak/Backend/internal/presentation/controller/v1/math"
	"github.com/Vanak/Backend/internal/presentation/controller/v1/user"
	"github.com/Vanak/Backend/internal/presentation/middleware"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type Application struct {
	Database              *Database
	MathController        *math.CustomerMathController
	Translator            domainLocalization.Translator
	Recovery              *middleware.RecoveryMiddleware
	Constants             *bootstrap.Constants
	OtpService            *service.OtpService
	UserService           *service.UserService
	GeneralUserController *user.GeneralUserController
}

func ProvideDBConfig(config *bootstrap.Config) *bootstrap.Database {
	return &config.Env.PrimaryDB
}

func ProvideConstants(cfg *bootstrap.Config) *bootstrap.Constants {
	return cfg.Constant
}

func ProvideRedisCache(cfg *bootstrap.Config) database.Cache {
	db := database.InitRedis(cfg)
	return &db
}

var MathProviderSet = wire.NewSet(
	math.NewCustomerMathController,
)

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDataBase,
	ProvideRedisCache,
	redisRepo.NewUserCacheRepository,
	wire.Bind(new(redisDomain.UserCacheRepository), new(*redisRepo.UserCacheRepository)),
	wire.Struct(new(Database), "*"),
)

var MiddlewareSet = wire.NewSet(
	middleware.NewRecoveryMiddleware,
)

var ServiceSet = wire.NewSet(
	mathService.NewMathService,
	localization.NewTranslationService,
	service.NewOtpService,
	service.NewUserService,
	service.NewJwtService,
	wire.Bind(new(usecase.MathService), new(*mathService.MathService)),
	wire.Bind(new(domainLocalization.Translator), new(*localization.TranslationService)),
	wire.Bind(new(usecase.OtpService), new(*service.OtpService)),
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
)

var ControllerSet = wire.NewSet(
	user.NewGeneralUserController,
)

func NewApplication(
	db *Database,
	mathCtrl *math.CustomerMathController,
	trans domainLocalization.Translator,
	rec *middleware.RecoveryMiddleware,
	cons *bootstrap.Constants,
	otpService *service.OtpService,
	userService *service.UserService,
	generalUserController *user.GeneralUserController,
) *Application {
	return &Application{
		Database:              db,
		MathController:        mathCtrl,
		Translator:            trans,
		Recovery:              rec,
		Constants:             cons,
		OtpService:            otpService,
		UserService:           userService,
		GeneralUserController: generalUserController,
	}
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	wire.Build(
		DatabaseProviderSet,
		MiddlewareSet,
		ServiceSet,
		MathProviderSet,
		ProvideDBConfig,
		ProvideConstants,
		NewApplication,
		ControllerSet,
	)
	return nil, nil
}
