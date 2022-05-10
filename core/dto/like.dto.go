package dto

import "github.com/google/uuid"

type LikedDto struct {
	User    BaseUserInfo `json:"me"`
	EventID *uuid.UUID   `json:"event_id"`
	PlaceID *uuid.UUID   `json:"place_id"`
}

type LikeDto struct {
	EventID *uuid.UUID `json:"event_id"`
	PlaceID *uuid.UUID `json:"place_id"`
}
