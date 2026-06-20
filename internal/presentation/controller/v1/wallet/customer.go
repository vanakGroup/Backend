package wallet

import (
	"fmt"

	"github.com/Vanak/Backend/bootstrap"
	walletdto "github.com/Vanak/Backend/internal/application/dto/wallet"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type CustomerWalletController struct {
	cfg           *bootstrap.Config
	walletService usecase.WalletService
}

func NewCustomerWalletController(cfg *bootstrap.Config, walletService usecase.WalletService) *CustomerWalletController {
	return &CustomerWalletController{
		cfg:           cfg,
		walletService: walletService,
	}
}

// GetWallet godoc
// @Summary      Get customer wallet details
// @Description  Retrieve the current balance and status of the authenticated customer's wallet.
// @Tags         Customer Wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  walletdto.WalletDTOResponse "Successful retrieval of wallet data"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      404  {object}  base.SingleMessageResponse "Wallet not found"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /wallet [get]
func (cwc *CustomerWalletController) GetWallet(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(cwc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cwc.cfg.Constant.ErrorField.User, cwc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userId := userIDAny.(uint)

	wallet, err := cwc.walletService.GetWalletInformation(userId)

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cwc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getWallet")
	base.Response(ctx, 200, message, wallet)
}

// Deposit godoc
// @Summary      Deposit funds into customer wallet
// @Description  Charge the authenticated customer's wallet using a mock gateway.
// @Tags         Customer Wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object{amount=integer} true "Deposit payload" example(object{amount=50000})
// @Success      200 {object} base.SingleMessageResponse "Successful registration message"
// @Failure      400 {object} base.SingleMessageResponse "Validation error"
// @Failure      401 {object} base.SingleMessageResponse "Unauthorized"
// @Failure      500 {object} base.SingleMessageResponse "Internal server error"
// @Router       /wallet/deposit [post]
func (cwc *CustomerWalletController) Deposit(ctx *gin.Context) {

	type WalletFilledParams struct {
		Amount int64 `json:"amount" validate:"required,gt=0"`
	}

	params := base.Validated[WalletFilledParams](ctx)

	userIDAny, exists := ctx.Get(cwc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cwc.cfg.Constant.ErrorField.User, cwc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userId := userIDAny.(uint)

	transactionInfo := walletdto.WalletDTORequest{
		Amount:          params.Amount,
		ReferenceType:   enums.GATEWAY,
		TransactionType: enums.DEPOSIT,
		UserId:          userId,
	}

	err := cwc.walletService.Deposit(transactionInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cwc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.registerTransaction")
	base.Response(ctx, 200, message, nil)

}

// GetLastTransactions godoc
// @Summary      Get recent wallet transactions
// @Description  Retrieve a list of the most recent transactions (deposits, withdrawals, etc.) for the authenticated customer's wallet.
// @Tags         Customer Wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   walletdto.TransactionListDTOResponse "List of recent transactions with details (amount, type, reason, timestamp)"
// @Failure      401  {object}  base.SingleMessageResponse "Unauthorized (missing or invalid token)"
// @Failure      500  {object}  base.SingleMessageResponse "Internal server error"
// @Router       /wallet/transactions [get]
func (cwc *CustomerWalletController) GetLastTransactions(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(cwc.cfg.Constant.Context.UserID)

	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cwc.cfg.Constant.ErrorField.User, cwc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userId := userIDAny.(uint)

	transactionList, err := cwc.walletService.GetLastTransactions(userId)

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cwc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getTransactionList")
	base.Response(ctx, 200, message, transactionList)
}
