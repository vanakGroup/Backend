package parkingdto

import "time"

type CreateParkingResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	AddressID    uint      `json:"address_id"`
	AttributeIDs []uint    `json:"attribute_ids,omitempty"`
	SlotCount    int       `json:"slot_count,omitempty"`
}
