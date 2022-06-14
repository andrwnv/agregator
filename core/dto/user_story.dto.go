package dto

import "github.com/google/uuid"

type UserStoryDto struct {
	ID           uuid.UUID        `json:"id"`
	Title        string           `json:"title"`
	LongReadText string           `json:"long_read_text"`
	CreatedBy    BaseUserInfo     `json:"created_by"`
	LinkedEvents []LinkedEventDto `json:"linked_events"`
	LinkedPlaces []LinkedPlaceDto `json:"linked_places"`
	LinkedPhotos []string         `json:"linked_photos"`
	CreatedAt    int64            `json:"created_at"`
}

type CreateUserStoryWithLinksDto struct {
	Title        string   `json:"title"`
	LongReadText string   `json:"long_read_text"`
	Events       []string `json:"linked_events"`
	Places       []string `json:"linked_places"`
}

type CreateUserStoryDto struct {
	Title        string `json:"title"`
	LongReadText string `json:"long_read_text"`
}

type UpdateUserStoryDto struct {
	Title         string   `json:"title"`
	LongReadText  string   `json:"long_read_text"`
	EventToDelete []string `json:"event_to_delete"`
	PlaceToDelete []string `json:"place_to_delete"`
	EventToCreate []string `json:"event_to_create"`
	PlaceToCreate []string `json:"place_to_create"`
}

type LinkedEventDto struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	StoryID string `json:"story_id"`
}

type LinkedPlaceDto struct {
	ID      string `json:"id"`
	PlaceID string `json:"place_id"`
	StoryID string `json:"story_id"`
}
