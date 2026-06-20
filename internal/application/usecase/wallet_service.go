package usecase

import (
	walletdto "github.com/Vanak/Backend/internal/application/dto/wallet"
)

type WalletService interface {
	GetWalletInformation(userId uint) (*walletdto.WalletDTOResponse, error)
	Deposit(transactionInfo walletdto.WalletDTORequest) error
	GetLastTransactions(userId uint) ([]walletdto.TransactionListDTOResponse, error)
}
