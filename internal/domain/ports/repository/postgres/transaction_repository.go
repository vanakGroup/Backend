package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type TransactionRepository interface {
	Create(*entities.Transaction) error
	FindLatestTransactions(walletId uint, n int) ([]entities.Transaction, error)
}
