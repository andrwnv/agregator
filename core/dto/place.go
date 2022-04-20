package dto

type CreatePlace struct {
	PaymentRequired bool    `json:"payment_required"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Longitude       float32 `json:"longitude"`
	Latitude        float32 `json:"latitude"`
}
