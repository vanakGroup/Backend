package persistance

import (
	"github.com/Vanak/Backend/bootstrap"
	postgresDomain "github.com/Vanak/Backend/internal/domain/ports/repository/postgres"
	"github.com/Vanak/Backend/internal/infrastructure/persistance/repository/postgres"
	"gorm.io/gorm"
)

type RepositoryFactory struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewRepositoryFactory(db *gorm.DB, cfg *bootstrap.Config) *RepositoryFactory {
	return &RepositoryFactory{
		db:  db,
		cfg: cfg,
	}
}

func (rf *RepositoryFactory) UserRepository() postgresDomain.UserRepository {
	return postgres.NewUserRepository(rf.cfg, rf.db)
}

func (rf *RepositoryFactory) AddressRepository() postgresDomain.AddressRepository {
	return postgres.NewAddressRepository(rf.cfg, rf.db)
}

func (rf *RepositoryFactory) ParkingRepository() postgresDomain.ParkingRepository {
	return postgres.NewParkingRepository(rf.cfg, rf.db)
}

func (f *RepositoryFactory) VehicleRepository() postgresDomain.VehicleRepository {
	return postgres.NewVehicleRepository(f.cfg, f.db)
}

func (f *RepositoryFactory) WalletRepository() postgresDomain.WalletRepository {
	return postgres.NewWalletRepository(f.cfg, f.db)
}

func (f *RepositoryFactory) PaymentRepository() postgresDomain.PaymentRepository {
	return postgres.NewPaymentRepository(f.cfg, f.db)
}

func (f *RepositoryFactory) TransactionRepository() postgresDomain.TransactionRepository {
	return postgres.NewTransactionRepository(f.cfg, f.db)
}
