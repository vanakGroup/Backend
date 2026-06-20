package entities

import (
	"github.com/Vanak/Backend/internal/domain/enums"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount        int64                          `gorm:"not null;type:bigint;default:0;check:amount >= 0"`
	Type          enums.TransactionType          `gorm:"not null;index"`
	ReferenceType enums.TransactionReferenceType `gorm:"not null;index:idx_tx_ref"`
	ReferenceID   uint                           `gorm:"not null;index:idx_tx_ref"`
	WalletID      uint                           `gorm:"not null;index"`
}
