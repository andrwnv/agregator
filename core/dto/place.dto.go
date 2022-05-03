package dto

import "github.com/google/uuid"

type PlaceDto struct {
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
	EventPhotos     []string     `json:"event_photos"`
}

type CreatePlace struct {
	PaymentRequired bool    `json:"payment_required" binding:"required"`
	Title           string  `json:"title" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	Longitude       float32 `json:"longitude" binding:"required"`
	Latitude        float32 `json:"latitude" binding:"required"`
	RegionID        string  `json:"region_id" binding:"required"`
}

type UpdatePlace struct {
	BeginDate       int64   `json:"begin_date"`
	EndDate         int64   `json:"end_date"`
	PaymentRequired bool    `json:"payment_required"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Longitude       float32 `json:"longitude"`
	Latitude        float32 `json:"latitude"`
	RegionID        string  `json:"region_id"`
}
