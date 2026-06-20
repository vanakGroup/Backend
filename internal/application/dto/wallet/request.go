package walletdto

import "github.com/Vanak/Backend/internal/domain/enums"

type WalletDTORequest struct {
	Amount          int64
	ReferenceType   enums.TransactionReferenceType
	TransactionType enums.TransactionType
	UserId          uint
}
