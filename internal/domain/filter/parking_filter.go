package filter

import "time"

type ParkingFilter struct {
	Price      PriceFilter
	Location   LocationFilter
	Attributes AttributeFilter
	Slot       *SlotFilter
}

type PriceFilter struct {
	Max *float64
	Min *float64
}

type LocationFilter struct {
	CityID  *uint
	Address *string
	GPS     *GPSFilter
}

type GPSFilter struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

type AttributeFilter struct {
	AttributeIDs []uint
}

type SlotFilter struct {
	Date      time.Time
	StartHour uint
	EndHour   uint
}
