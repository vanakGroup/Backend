package enums

type TransactionType uint

const (
	DEPOSIT TransactionType = iota + 1
	PAYMENT
	REFUND
)

func (tt TransactionType) String() string {
	switch tt {
	case DEPOSIT:
		return "deposit"
	case PAYMENT:
		return "payment"
	case REFUND:
		return "refund"
	}

	return ""
}

func GetAllTransactionTypes() []TransactionType {
	return []TransactionType{
		DEPOSIT,
		PAYMENT,
		REFUND,
	}
}
