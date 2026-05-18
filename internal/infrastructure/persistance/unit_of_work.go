package persistance

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/ports"
	"gorm.io/gorm"
)

type UnitOfWork struct {
	db  *gorm.DB
	cfg *bootstrap.Config
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (u *UnitOfWork) Factory() ports.RepositoryFactory {
	return NewRepositoryFactory(u.db, u.cfg)
}

func (u *UnitOfWork) WithTransaction(fn func(ports.RepositoryFactory) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		factory := NewRepositoryFactory(tx, u.cfg)
		return fn(factory)
	})
}
