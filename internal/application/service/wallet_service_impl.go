package service

import (
	"github.com/Vanak/Backend/bootstrap"
	walletdto "github.com/Vanak/Backend/internal/application/dto/wallet"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/infrastructure/hashed"
)

type WalletService struct {
	unitOfWork ports.UnitOfWork
	cfg        *bootstrap.Config
	hashID     *hashed.HashID
}

func NewWalletService(
	unitOfWork ports.UnitOfWork,
	cfg *bootstrap.Config,
	hashID *hashed.HashID,
) *WalletService {
	return &WalletService{
		unitOfWork: unitOfWork,
		cfg:        cfg,
		hashID:     hashID,
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

func (ws *WalletService) depositInternal(rf ports.RepositoryFactory, transactionInfo walletdto.WalletDTORequest) error {
	foundWallet, err := ws.unitOfWork.Factory().WalletRepository().FindWalletByUserId(transactionInfo.UserId)
	disableErrors := exception.NewDisabledError()

	if err != nil {
		return err
	}

	if !foundWallet.IsActive {
		disableErrors.Add(ws.cfg.Constant.ErrorField.Wallet, ws.cfg.Constant.ErrorTag.Disabled)
	}

	if transactionInfo.TransactionType != enums.DEPOSIT {
		return exception.NewConflictError(ws.cfg.Constant.ErrorField.Transaction)
	}

	if disableErrors.HasErrors() {
		return disableErrors
	}

	paymentRepo := rf.PaymentRepository()
	transactionRepo := rf.TransactionRepository()
	walletRepo := rf.WalletRepository()

	lockedWallet, err := walletRepo.FindByIdForUpdate(foundWallet.ID)
	if err != nil {
		return err
	}

	var referenceId uint
	if transactionInfo.ReferenceType == enums.GATEWAY {
		payment := &entities.MockPayment{
			Amount: transactionInfo.Amount,
			UserID: transactionInfo.UserId,
		}

		err = paymentRepo.Create(payment)
		if err != nil {
			return err
		}

		referenceId = payment.ID
	} else if transactionInfo.ReferenceType == enums.RESERVATION {
		if transactionInfo.ReferenceId == nil {
			return exception.NewConflictError(ws.cfg.Constant.ErrorField.Transaction)
		}
		referenceId = *transactionInfo.ReferenceId
	} else {
		return exception.NewConflictError(ws.cfg.Constant.ErrorField.Transaction)
	}

	transaction := &entities.Transaction{
		Amount:        transactionInfo.Amount,
		Type:          transactionInfo.TransactionType,
		ReferenceType: transactionInfo.ReferenceType,
		ReferenceID:   referenceId,
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

	return nil
}

func (ws *WalletService) deductInternal(rf ports.RepositoryFactory, transactionInfo walletdto.WalletDTORequest) error {
	foundWallet, err := ws.unitOfWork.Factory().WalletRepository().FindWalletByUserId(transactionInfo.UserId)
	disableErrors := exception.NewDisabledError()

	if err != nil {
		return err
	}

	if transactionInfo.TransactionType != enums.PAYMENT {
		return exception.NewConflictError(ws.cfg.Constant.ErrorField.Transaction)
	}

	if !foundWallet.IsActive {
		disableErrors.Add(ws.cfg.Constant.ErrorField.Wallet, ws.cfg.Constant.ErrorTag.Disabled)
	}

	if disableErrors.HasErrors() {
		return disableErrors
	}

	transactionRepo := rf.TransactionRepository()
	walletRepo := rf.WalletRepository()

	lockedWallet, err := walletRepo.FindByIdForUpdate(foundWallet.ID)
	if err != nil {
		return err
	}

	if lockedWallet.Balance < transactionInfo.Amount {
		return exception.NewConflictError(ws.cfg.Constant.ErrorField.Balance)
	}

	transaction := &entities.Transaction{
		Amount:        transactionInfo.Amount,
		Type:          transactionInfo.TransactionType,
		ReferenceType: transactionInfo.ReferenceType,
		ReferenceID:   *transactionInfo.ReferenceId,
		WalletID:      foundWallet.ID,
	}

	err = transactionRepo.Create(transaction)
	if err != nil {
		return err
	}

	err = walletRepo.UpdateBalance(lockedWallet.ID, lockedWallet.Balance-transaction.Amount)
	if err != nil {
		return err
	}

	return nil

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
		transactionID, err := ws.hashID.Encode(t.ID)
		if err != nil {
			return nil, err
		}
		transactionList = append(transactionList, walletdto.TransactionListDTOResponse{
			ID:        transactionID,
			Amount:    t.Amount,
			Type:      t.Type.String(),
			Reason:    t.ReferenceType.String(),
			CreatedAt: t.CreatedAt.Unix(),
		})
	}

	return transactionList, nil
}

func (ws *WalletService) Deposit(transactionInfo walletdto.WalletDTORequest) error {
	return ws.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		return ws.depositInternal(rf, transactionInfo)
	})
}

func (ws *WalletService) Deduct(transactionInfo walletdto.WalletDTORequest) error {
	return ws.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		return ws.deductInternal(rf, transactionInfo)
	})
}

func (ws *WalletService) Transfer(deductReq walletdto.WalletDTORequest, depositReq walletdto.WalletDTORequest) error {
	return ws.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {

		if err := ws.deductInternal(rf, deductReq); err != nil {
			return err
		}

		if err := ws.depositInternal(rf, depositReq); err != nil {
			return err
		}

		return nil
	})
}
