package dto

import "github.com/google/uuid"

type PlaceDto struct {
	ID          uuid.UUID    `json:"id"`
	PaymentNeed bool         `json:"payment_need"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Longitude   float32      `json:"longitude"`
	Latitude    float32      `json:"latitude"`
	CreatedBy   BaseUserInfo `json:"created_by"`
	RegionInfo  RegionDto    `json:"region_info"`
	PlacePhotos []string     `json:"photos"`
}

type CreatePlace struct {
	PaymentNeed bool    `json:"payment_need"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Longitude   float32 `json:"longitude" binding:"required"`
	Latitude    float32 `json:"latitude" binding:"required"`
	RegionID    string  `json:"region_id" binding:"required"`
}

type UpdatePlace struct {
	PaymentNeed bool    `json:"payment_need"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Longitude   float32 `json:"longitude"`
	Latitude    float32 `json:"latitude"`
	RegionID    string  `json:"region_id"`
}
