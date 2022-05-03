package dto

// ----- Event comments -----

type EventCommentDto struct {
	ID            string       `json:"id"`
	CreatedBy     BaseUserInfo `json:"created_by"`
	LinkedEventID string       `json:"linked_event_id"`
	CommentBody   string       `json:"comment_body"`
	UpdatedAt     int64        `json:"updated_at"`
	CreatedAt     int64        `json:"created_at"`
}

type CreateEventCommentDto struct {
	LinkedEventID string `json:"linked_event_id" binding:"required"`
	CommentBody   string `json:"comment_body" binding:"required"`
}

type UpdateEventCommentDto struct {
	CommentBody string `json:"comment_body" binding:"required"`
}

// ----- Place comments -----

type PlaceCommentDto struct {
	ID            string       `json:"id"`
	CreatedBy     BaseUserInfo `json:"created_by"`
	LinkedPlaceID string       `json:"linked_place_id"`
	CommentBody   string       `json:"comment_body"`
	UpdatedAt     int64        `json:"updated_at"`
	CreatedAt     int64        `json:"created_at"`
}

type CreatePlaceCommentDto struct {
	LinkedPlaceID string `json:"linked_place_id" binding:"required"`
	CommentBody   string `json:"comment_body" binding:"required"`
}

type UpdatePlaceCommentDto struct {
	CommentBody string `json:"comment_body" binding:"required"`
}
