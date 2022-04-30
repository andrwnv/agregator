package dto

import "github.com/google/uuid"

type EventDto struct {
	ID              uuid.UUID    `json:"id"`
	BeginDate       int64        `json:"begin_date"`
	EndDate         int64        `json:"end_date"`
	PaymentRequired bool         `json:"payment_required"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Longitude       float32      `json:"longitude"`
	Latitude        float32      `json:"latitude"`
	CreatedBy       BaseUserInfo `json:"created_by"`
	RegionInfo      RegionDto    `json:"region_info"`
}

type CreateEvent struct {
	BeginDate       int64   `json:"begin_date"`
	EndDate         int64   `json:"end_date"`
	PaymentRequired bool    `json:"payment_required"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Longitude       float32 `json:"longitude"`
	Latitude        float32 `json:"latitude"`
	RegionID        string  `json:"region_id"`
}
