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
	"github.com/Vanak/Backend/internal/presentation/middleware"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type Application struct {
	Database       *Database
	MathController *math.CustomerMathController
	Translator     domainLocalization.Translator
	Recovery       *middleware.RecoveryMiddleware
	Constants      *bootstrap.Constants
	OtpService     *service.OtpService
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
	wire.Bind(new(usecase.MathService), new(*mathService.MathService)),
	wire.Bind(new(domainLocalization.Translator), new(*localization.TranslationService)),
)

func NewApplication(
	db *Database,
	mathCtrl *math.CustomerMathController,
	trans domainLocalization.Translator,
	rec *middleware.RecoveryMiddleware,
	cons *bootstrap.Constants,
	otpService *service.OtpService,
) *Application {
	return &Application{
		Database:       db,
		MathController: mathCtrl,
		Translator:     trans,
		Recovery:       rec,
		Constants:      cons,
		OtpService:     otpService,
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
	)
	return nil, nil
}
