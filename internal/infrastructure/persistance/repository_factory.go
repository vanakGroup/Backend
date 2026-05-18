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

func (f *RepositoryFactory) UserRepository() postgresDomain.UserRepository {
	return postgres.NewUserRepository(f.cfg, f.db)
}
