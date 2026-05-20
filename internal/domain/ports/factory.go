package ports

import postgresDomain "github.com/Vanak/Backend/internal/domain/ports/repository/postgres"

type RepositoryFactory interface {
	UserRepository() postgresDomain.UserRepository
	AddressRepository() postgresDomain.AddressRepository
}
