package enums

type TransactionReferenceType uint

const (
	RESERVATION TransactionReferenceType = iota + 1
	GATEWAY
)

func (trt TransactionReferenceType) String() string {
	switch trt {
	case RESERVATION:
		return "reservation"
	case GATEWAY:
		return "gateway"
	}

	return ""
}

func GetAllTransactionRefrenceTypes() []TransactionReferenceType {
	return []TransactionReferenceType{
		RESERVATION,
		GATEWAY,
	}
}
