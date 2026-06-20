package enums

type ReservationStatus uint

const (
	PENDING ReservationStatus = iota + 1
	CONFIRMED
	COMPLETED
	CANCELLED
)

func (rs ReservationStatus) String() string {
	switch rs {
	case PENDING:
		return "pending"
	case CONFIRMED:
		return "confirmed"
	case COMPLETED:
		return "completed"
	case CANCELLED:
		return "cancelled"
	}

	return ""
}

func GetAllReservationsStatus() []ReservationStatus {
	return []ReservationStatus{
		PENDING,
		CONFIRMED,
		COMPLETED,
		CANCELLED,
	}
}
