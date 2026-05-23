package enums

type BucketType uint

const (
	LogoPic BucketType = iota + 1
	ParkingDoc
	ParkingPic
)

func (bt BucketType) String() string {
	switch bt {
	case LogoPic:
		return "logoPic"
	case ParkingDoc:
		return "parkingDoc"
	case ParkingPic:
		return "parkingPic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		LogoPic,
		ParkingDoc,
		ParkingPic,
	}
}
