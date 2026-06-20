package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewPaymentRepository(cfg *bootstrap.Config, db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		cfg: cfg,
		db:  db,
	}
}

func (pr *PaymentRepository) Create(payment *entities.MockPayment) error {
	err := pr.db.Create(payment).Error

	if err != nil {
		return err
	}

	return nil
}
