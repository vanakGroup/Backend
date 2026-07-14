package walletdto

type WalletDTOResponse struct {
	Balance  int64 `json:"balance"`
	IsActive bool  `json:"is_active"`
}

type TransactionListDTOResponse struct {
	ID        string `json:"id"`
	Amount    int64  `json:"amount"`
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	CreatedAt int64  `json:"created_at"`
}
