package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewTransactionRepository(cfg *bootstrap.Config, db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		cfg: cfg,
		db:  db,
	}
}

func (tr *TransactionRepository) Create(transaction *entities.Transaction) error {
	err := tr.db.Create(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

func (tr *TransactionRepository) FindLatestTransactions(walledId uint, n int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	err := tr.db.Where("wallet_id = ?", walledId).
		Order("created_at desc").Limit(n).Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
