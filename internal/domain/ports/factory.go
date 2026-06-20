package ports

import postgresDomain "github.com/Vanak/Backend/internal/domain/ports/repository/postgres"

type RepositoryFactory interface {
	UserRepository() postgresDomain.UserRepository
	AddressRepository() postgresDomain.AddressRepository
	ParkingRepository() postgresDomain.ParkingRepository
	VehicleRepository() postgresDomain.VehicleRepository
	WalletRepository() postgresDomain.WalletRepository
	PaymentRepository() postgresDomain.PaymentRepository
	TransactionRepository() postgresDomain.TransactionRepository
}
