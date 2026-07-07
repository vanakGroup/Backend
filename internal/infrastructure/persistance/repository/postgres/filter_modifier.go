package postgres

import (
	"github.com/Vanak/Backend/internal/domain/filter"
	"gorm.io/gorm"
)

type ParkingQueryBuilder struct {
	query *gorm.DB
}

func NewParkingQueryBuilder(query *gorm.DB) *ParkingQueryBuilder {
	return &ParkingQueryBuilder{query: query}
}

func (pqb *ParkingQueryBuilder) ApplyPrice(price filter.PriceFilter) *ParkingQueryBuilder {
	query := pqb.query

	if price.Min != nil {
		query = query.Where("base_price >= ?", *price.Min)
	}

	if price.Max != nil {
		query = query.Where("base_price <= ?", *price.Max)
	}

	pqb.query = query

	return pqb
}

func (pqb *ParkingQueryBuilder) ApplyLocation(location filter.LocationFilter) *ParkingQueryBuilder {
	query := pqb.query

	if location.CityID == nil &&
		location.Address == nil &&
		location.GPS == nil {
		return pqb
	}

	query = query.Joins(`
		LEFT JOIN addresses
			ON addresses.parking_id = parkings.id
	`)

	if location.CityID != nil {
		query = query.Where("addresses.city_id = ?", *location.CityID)
	}

	if location.Address != nil {
		query = query.Where(
			"addresses.exact_address ILIKE ?",
			"%"+*location.Address+"%",
		)
	}

	if location.GPS != nil {
		query = query.Where(`
		6371 * ACOS(
			COS(RADIANS(?)) *
			COS(RADIANS(parkings.latitude)) *
			COS(RADIANS(parkings.longitude) - RADIANS(?)) +
			SIN(RADIANS(?)) *
			SIN(RADIANS(parkings.latitude))
		) <= ?
	`,
			location.GPS.Latitude,
			location.GPS.Longitude,
			location.GPS.Latitude,
			location.GPS.Radius,
		)
	}

	pqb.query = query

	return pqb
}

func (pqb *ParkingQueryBuilder) ApplyAttributes(attr filter.AttributeFilter) *ParkingQueryBuilder {
	query := pqb.query

	if len(attr.AttributeIDs) > 0 {
		query = query.
			Joins(`
				JOIN parking_attribute_relations par
					ON par.parking_id = parkings.id
			`).
			Where(
				"par.attribute_id IN ?",
				attr.AttributeIDs,
			).
			Group("parkings.id").
			Having(
				"COUNT(DISTINCT par.attribute_id) = ?",
				len(attr.AttributeIDs),
			)
	}

	pqb.query = query

	return pqb
}

func (pqb *ParkingQueryBuilder) ApplySlot(slot *filter.SlotFilter) *ParkingQueryBuilder {
	query := pqb.query

	if slot != nil {
		requestedHours := slot.EndHour - slot.StartHour

		query = query.
			Joins(`
			JOIN availability_slots availability
				ON availability.parking_id = parkings.id
		`).
			Where("availability.date = ?", slot.Date).
			Where("availability.hour >= ?", slot.StartHour).
			Where("availability.hour < ?", slot.EndHour).
			Where("availability.is_booked = ?", false).
			Group("parkings.id").
			Having("COUNT(DISTINCT availability.hour) = ?", requestedHours)
	}

	pqb.query = query

	return pqb
}

func (pqb *ParkingQueryBuilder) Apply(filter filter.ParkingFilter) *ParkingQueryBuilder {
	return pqb.
		ApplyPrice(filter.Price).
		ApplyLocation(filter.Location).
		ApplyAttributes(filter.Attributes).
		ApplySlot(filter.Slot)
}

func (pqb *ParkingQueryBuilder) Build() *gorm.DB {
	return pqb.query
}
