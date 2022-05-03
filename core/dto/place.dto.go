package dto

import "github.com/google/uuid"

type PlaceDto struct {
	ID              uuid.UUID    `json:"id"`
	PaymentRequired bool         `json:"payment_required"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Longitude       float32      `json:"longitude"`
	Latitude        float32      `json:"latitude"`
	CreatedBy       BaseUserInfo `json:"created_by"`
	RegionInfo      RegionDto    `json:"region_info"`
	PlacePhotos     []string     `json:"place_photos"`
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
	PaymentRequired bool    `json:"payment_required"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Longitude       float32 `json:"longitude"`
	Latitude        float32 `json:"latitude"`
	RegionID        string  `json:"region_id"`
}
