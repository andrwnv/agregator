package dto

import "time"

type CreateEvent struct {
	BeginDate       time.Time `json:"begin_date"`
	EndDate         time.Time `json:"end_date"`
	PaymentRequired bool      `json:"payment_required"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Longitude       float32   `json:"longitude"`
	Latitude        float32   `json:"latitude"`
}
