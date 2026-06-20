package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type WalletRepository interface {
	Create(*entities.Wallet) error
	UpdateBalance(walletId uint, balance int64) error
	ActiveWallet(userId uint, value bool) error
	FindWalletByUserId(userId uint) (*entities.Wallet, error)
	FindByIdForUpdate(walledId uint) (*entities.Wallet, error)
}
