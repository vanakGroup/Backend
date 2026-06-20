package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewWalletRepository(cfg *bootstrap.Config, db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		cfg: cfg,
		db:  db,
	}
}

func (wr *WalletRepository) Create(Wallet *entities.Wallet) error {
	return wr.db.Create(Wallet).Error
}

func (wr *WalletRepository) ActiveWallet(userId uint, value bool) error {
	err := wr.db.Model(&entities.Wallet{}).
		Where("user_id = ?", userId).
		Update("IsActive", value).Error

	if err != nil {
		return err
	}

	return nil
}

func (wr *WalletRepository) FindWalletByUserId(userId uint) (*entities.Wallet, error) {
	var wallet entities.Wallet

	err := wr.db.Model(&entities.Wallet{}).
		Where("user_id = ?", userId).
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (wr *WalletRepository) FindByIdForUpdate(walletId uint) (*entities.Wallet, error) {
	var wallet entities.Wallet

	err := wr.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", walletId).
		First(&wallet).Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (wr *WalletRepository) UpdateBalance(walletId uint, balance int64) error {
	err := wr.db.Model(&entities.Wallet{}).
		Where("id = ?", walletId).
		Update("Balance", balance).Error

	if err != nil {
		return err
	}

	return nil
}
