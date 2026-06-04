package parkingdto

import "time"

type CreateParkingResponse struct {
	ID           uint      `json:"id" example:"1"`
	Title        string    `json:"title" example:"پارکینگ طبقاتی ونک"`
	Status       string    `json:"status" example:"Active"`
	CreatedAt    time.Time `json:"createdAt" example:"2026-06-04T00:00:00Z"`
	AddressID    uint      `json:"addressId" example:"12"`
	AttributeIDs []uint    `json:"attributeIds,omitempty" example:"[1,2,3]"`
	SlotCount    int       `json:"slotCount" example:"24"`
	
	DocumentURL  string    `json:"documentUrl" example:"https://s3.ir-thr-at1.arvanstorage.ir/vanak/parkings/1/document.pdf"`
	ImageURLs    []string  `json:"imageUrls,omitempty" example:"[\"https://s3.ir-thr-at1.arvanstorage.ir/vanak/parkings/1/images/image_0.jpg\",\"https://s3.ir-thr-at1.arvanstorage.ir/vanak/parkings/1/images/image_1.jpg\"]"`
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