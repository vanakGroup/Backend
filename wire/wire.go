//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/application/service"
	"github.com/Vanak/Backend/internal/application/usecase"
	jwtDomain "github.com/Vanak/Backend/internal/domain/jwt"
	domainLocalization "github.com/Vanak/Backend/internal/domain/localization"
	"github.com/Vanak/Backend/internal/domain/ports"
	postgresDomain "github.com/Vanak/Backend/internal/domain/ports/repository/postgres"
	redisDomain "github.com/Vanak/Backend/internal/domain/ports/repository/redis"
	storageDomain "github.com/Vanak/Backend/internal/domain/storage"
	"github.com/Vanak/Backend/internal/infrastructure/database"
	"github.com/Vanak/Backend/internal/infrastructure/jwt"
	"github.com/Vanak/Backend/internal/infrastructure/localization"
	"github.com/Vanak/Backend/internal/infrastructure/persistance"
	postgresRepo "github.com/Vanak/Backend/internal/infrastructure/persistance/repository/postgres"
	redisRepo "github.com/Vanak/Backend/internal/infrastructure/persistance/repository/redis"
	"github.com/Vanak/Backend/internal/infrastructure/seeder"
	storage "github.com/Vanak/Backend/internal/infrastructure/storage"
	"github.com/Vanak/Backend/internal/presentation/controller/v1/parking"
	"github.com/Vanak/Backend/internal/presentation/controller/v1/user"
	"github.com/Vanak/Backend/internal/presentation/controller/v1/vehicle"
	"github.com/Vanak/Backend/internal/presentation/middleware"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type Middlewares struct {
	Authentication *middleware.AuthMiddleware
	Recovery       *middleware.RecoveryMiddleware
}

type Seeds struct {
	Address   *seeder.AddressSeeder
	Attribute *seeder.AttributeTypeSeeder
}

type Application struct {
	Database                  *Database
	Translator                domainLocalization.Translator
	Middlewares               *Middlewares
	Constants                 *bootstrap.Constants
	OtpService                *service.OtpService
	UserService               *service.UserService
	GeneralUserController     *user.GeneralUserController
	Seeds                     *Seeds
	CustomerParkingController *parking.CustomerParkingController
	CustomerUserController    *user.CustomerUserController
	CustomerVehicleController *vehicle.CustomerVehicleController
}

func ProvideDBConfig(config *bootstrap.Config) *bootstrap.Database {
	return &config.Env.PrimaryDB
}

// func ProvideS3Config(config *bootstrap.Config) *bootstrap.S3 {
// 	return &config.Env.Storage
// }

func ProvideConstants(cfg *bootstrap.Config) *bootstrap.Constants {
	return cfg.Constant
}

func ProvideRedisCache(cfg *bootstrap.Config) database.Cache {
	db := database.InitRedis(cfg)
	return &db
}

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDataBase,
	ProvideRedisCache,
	redisRepo.NewUserCacheRepository,
	postgresRepo.NewUserRepository,
	postgresRepo.NewAddressRepository,
	postgresRepo.NewParkingRepository,
	postgresRepo.NewVehicleRepository,
	persistance.NewUnitOfWork,
	persistance.NewRepositoryFactory,
	storage.NewS3Storage,
	wire.Bind(new(storageDomain.S3Storage), new(*storage.S3Storage)),
	wire.Bind(new(ports.UnitOfWork), new(*persistance.UnitOfWork)),
	wire.Bind(new(ports.RepositoryFactory), new(*persistance.RepositoryFactory)),
	wire.Bind(new(redisDomain.UserCacheRepository), new(*redisRepo.UserCacheRepository)),
	wire.Bind(new(postgresDomain.UserRepository), new(*postgresRepo.UserRepository)),
	wire.Bind(new(postgresDomain.AddressRepository), new(*postgresRepo.AddressRepository)),
	wire.Bind(new(postgresDomain.ParkingRepository), new(*postgresRepo.ParkingRepository)),
	wire.Bind(new(postgresDomain.UserProfileRepository), new(*postgresRepo.UserRepository)),
	wire.Bind(new(postgresDomain.VehicleRepository), new(*postgresRepo.VehicleRepository)),
	wire.Struct(new(Database), "*"),
)

var MiddlewareSet = wire.NewSet(
	middleware.NewRecoveryMiddleware,
	middleware.NewAuthMiddleWare,
	wire.Struct(new(Middlewares), "*"),
)

var ServiceSet = wire.NewSet(
	localization.NewTranslationService,
	service.NewOtpService,
	service.NewUserService,
	service.NewJwtService,
	service.NewParkingService,
	jwt.NewJwtKeyManager,
	service.NewVehicleService,
	wire.Bind(new(usecase.UserProfileUseCase), new(*service.UserService)),
	wire.Bind(new(usecase.VehicleService), new(*service.VehicleService)),
	wire.Bind(new(domainLocalization.Translator), new(*localization.TranslationService)),
	wire.Bind(new(jwtDomain.KeyManager), new(*jwt.JWTKeyManager)),
	wire.Bind(new(usecase.OtpService), new(*service.OtpService)),
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.JwtService), new(*service.JwtService)),
	wire.Bind(new(usecase.ParkingService), new(*service.ParkingService)),
)

var SeederSet = wire.NewSet(
	seeder.NewAddressSeeder,
	seeder.NewAttributeTypeSeeder,
	wire.Struct(new(Seeds), "*"),
)

var ControllerSet = wire.NewSet(
	user.NewGeneralUserController,
	parking.NewCustomerParkingController,
	user.NewCustomerUserController,
	vehicle.NewVehicleCustomerController,
)

func NewApplication(
	db *Database,
	trans domainLocalization.Translator,
	mid *Middlewares,
	cons *bootstrap.Constants,
	otpService *service.OtpService,
	userService *service.UserService,
	generalUserController *user.GeneralUserController,
	customerParkingController *parking.CustomerParkingController,
	seeds *Seeds,
	customerUserController *user.CustomerUserController,
	customerVehicleController *vehicle.CustomerVehicleController,
) *Application {
	return &Application{
		Database:                  db,
		Translator:                trans,
		Middlewares:               mid,
		Constants:                 cons,
		OtpService:                otpService,
		UserService:               userService,
		GeneralUserController:     generalUserController,
		Seeds:                     seeds,
		CustomerParkingController: customerParkingController,
		CustomerUserController:    customerUserController,
		CustomerVehicleController: customerVehicleController,
	}
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	wire.Build(
		DatabaseProviderSet,
		MiddlewareSet,
		ServiceSet,
		ProvideDBConfig,
		// ProvideS3Config,
		ProvideConstants,
		NewApplication,
		ControllerSet,
		SeederSet,
	)
	return nil, nil
}
