package parkingdto

import "time"

type CreateParkingResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	AddressID    uint      `json:"addressId"`
	AttributeIDs []uint    `json:"attributeIds,omitempty"`
	SlotCount    int       `json:"slotCount"`

	DocumentURL string   `json:"documentUrl"`
	ImageURLs   []string `json:"imageUrls,omitempty"`
}

type ParkingResponse struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	BasePrice    uint     `json:"basePrice"`
	DocumentURL  string   `json:"documentUrl,omitempty"`
	ImageURLs    []string `json:"imageUrls,omitempty"`
	PostalCode   string   `json:"postalCode"`
	ExactAddress string   `json:"exactAddress"`
}

type ParkingDetailResponse struct {
	ID           string                `json:"id"`
	Title        string                `json:"title"`
	Latitude     float64               `json:"latitude"`
	Longitude    float64               `json:"longitude"`
	BasePrice    uint                  `json:"basePrice"`
	Status       string                `json:"status"`
	DocumentURL  string                `json:"documentUrl"`
	ImageURLs    []string              `json:"imageUrls"`
	Address      AddressDetailResponse `json:"address"`
	Slots        []SlotDetailResponse  `json:"slots"`
	AttributeIDs []uint                `json:"attributeIds"`
}

type AddressDetailResponse struct {
	CityID       uint   `json:"cityId"`
	PostalCode   string `json:"postalCode"`
	UnitNumber   uint   `json:"unitNumber"`
	ExactAddress string `json:"exactAddress"`
}

type SlotDetailResponse struct {
	ID       uint      `json:"id"`
	Date     time.Time `json:"date"`
	Hour     uint      `json:"hour"`
	IsBooked bool      `json:"isBooked"`
}

type ParkingInfoResponse struct {
	ID           string                `json:"id"`
	Title        string                `json:"title"`
	Latitude     float64               `json:"latitude"`
	Longitude    float64               `json:"longitude"`
	BasePrice    float64               `json:"base_price"`
	Status       string                `json:"status"`
	Address      AddressDetailResponse `json:"address"`
	AttributeIDs []uint                `json:"attribute_ids"`
	ImageURLs    []string              `json:"image_urls"`
}

type ParkingForListResponse struct {
	ID           string                `json:"id"`
	Title        string                `json:"title"`
	BasePrice    float64               `json:"base_price"`
	Status       string                `json:"status"`
	Address      AddressDetailResponse `json:"address"`
	AttributeIDs []uint                `json:"attribute_ids"`
	ImageURLs    []string              `json:"image_urls"`
}

type ParkingReservationDetailResponse struct {
	ID        string                `json:"id"`
	Title     string                `json:"title"`
	Latitude  float64               `json:"latitude"`
	Longitude float64               `json:"longitude"`
	Address   AddressDetailResponse `json:"address"`
}
type ParkingSearchPaginationMeta struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

type ParkingSearchResponse struct {
	Data       []ParkingForListResponse    `json:"data"`
	Pagination ParkingSearchPaginationMeta `json:"pagination"`
}
