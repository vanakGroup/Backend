package service

import (
	"github.com/Vanak/Backend/bootstrap"
	walletdto "github.com/Vanak/Backend/internal/application/dto/wallet"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/ports"
)

type WalletService struct {
	unitOfWork ports.UnitOfWork
	cfg        *bootstrap.Config
}

func NewWalletService(
	unitOfWork ports.UnitOfWork,
	cfg *bootstrap.Config,
) *WalletService {
	return &WalletService{
		unitOfWork: unitOfWork,
		cfg:        cfg,
	}
}

func (ws *WalletService) GetWalletInformation(userId uint) (*walletdto.WalletDTOResponse, error) {
	foundWallet, err := ws.unitOfWork.Factory().WalletRepository().FindWalletByUserId(userId)
	disableErrors := exception.NewDisabledError()

	if err != nil {
		return nil, err
	}

	if !foundWallet.IsActive {
		disableErrors.Add(ws.cfg.Constant.ErrorField.Wallet, ws.cfg.Constant.ErrorTag.Disabled)
	}

	if disableErrors.HasErrors() {
		return nil, disableErrors
	}

	wallet := walletdto.WalletDTOResponse{
		Balance:  foundWallet.Balance,
		IsActive: foundWallet.IsActive,
	}

	return &wallet, nil
}

func (ws *WalletService) Deposit(transactionInfo walletdto.WalletDTORequest) error {
	foundWallet, err := ws.unitOfWork.Factory().WalletRepository().FindWalletByUserId(transactionInfo.UserId)
	disableErrors := exception.NewDisabledError()

	if err != nil {
		return err
	}

	if !foundWallet.IsActive {
		disableErrors.Add(ws.cfg.Constant.ErrorField.Wallet, ws.cfg.Constant.ErrorTag.Disabled)
	}

	if disableErrors.HasErrors() {
		return disableErrors
	}

	err = ws.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		paymentRepo := rf.PaymentRepository()
		transactionRepo := rf.TransactionRepository()
		walletRepo := rf.WalletRepository()

		if transactionInfo.TransactionType == enums.DEPOSIT && transactionInfo.ReferenceType == enums.GATEWAY {
			lockedWallet, err := walletRepo.FindByIdForUpdate(foundWallet.ID)
			if err != nil {
				return err
			}

			payment := &entities.MockPayment{
				Amount: transactionInfo.Amount,
				UserID: transactionInfo.UserId,
			}

			err = paymentRepo.Create(payment)
			if err != nil {
				return err
			}

			transaction := &entities.Transaction{
				Amount:        transactionInfo.Amount,
				Type:          transactionInfo.TransactionType,
				ReferenceType: transactionInfo.ReferenceType,
				ReferenceID:   payment.ID,
				WalletID:      foundWallet.ID,
			}

			err = transactionRepo.Create(transaction)
			if err != nil {
				return err
			}

			err = walletRepo.UpdateBalance(lockedWallet.ID, lockedWallet.Balance+transaction.Amount)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err

}

func (ws *WalletService) GetLastTransactions(userId uint) ([]walletdto.TransactionListDTOResponse, error) {
	foundWallet, err := ws.unitOfWork.Factory().WalletRepository().FindWalletByUserId(userId)
	disableErrors := exception.NewDisabledError()

	if err != nil {
		return nil, err
	}

	if !foundWallet.IsActive {
		disableErrors.Add(ws.cfg.Constant.ErrorField.Wallet, ws.cfg.Constant.ErrorTag.Disabled)
	}

	if disableErrors.HasErrors() {
		return nil, disableErrors
	}

	transactions, err := ws.unitOfWork.Factory().TransactionRepository().FindLatestTransactions(foundWallet.ID, 10)

	if err != nil {
		return nil, err
	}

	transactionList := make([]walletdto.TransactionListDTOResponse, 0)
	for _, t := range transactions {
		transactionList = append(transactionList, walletdto.TransactionListDTOResponse{
			ID:        t.ID,
			Amount:    t.Amount,
			Type:      t.Type.String(),
			Reason:    t.ReferenceType.String(),
			CreatedAt: t.CreatedAt.Unix(),
		})
	}

	return transactionList, nil
}
