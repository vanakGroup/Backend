package enums

type ParkingStatus uint

const (
	ParkingStatusActive ParkingStatus = iota + 1
	ParkingStatusInActive
)

func (ps ParkingStatus) String() string {
	switch ps {
	case ParkingStatusActive:
		return "فعال"
	case ParkingStatusInActive:
		return "غیرفعال"
	}
	return "نامعلوم"
}

func GetAllParkingStatus() []ParkingStatus {
	return []ParkingStatus{
		ParkingStatusActive,
		ParkingStatusInActive,
	}
}
