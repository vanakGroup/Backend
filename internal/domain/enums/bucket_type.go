package enums

type BucketType uint

const (
	LogoPic BucketType = iota + 1
	ParkingDoc
	ParkingPic
	RegistrationCardPic
)

func (bt BucketType) String() string {
	switch bt {
	case LogoPic:
		return "logoPic"
	case ParkingDoc:
		return "parkingDoc"
	case ParkingPic:
		return "parkingPic"
	case RegistrationCardPic:
		return "registrationCarPic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		LogoPic,
		ParkingDoc,
		ParkingPic,
		RegistrationCardPic,
	}
}
