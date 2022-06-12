package repo

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model

	ID              uuid.UUID `gorm:"primaryKey"`
	BeginDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	PaymentRequired bool      `gorm:"default: false"`
	CreatedByID     uuid.UUID `gorm:"not null"`
	Title           string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	Longitude       float32   `gorm:"not null"`
	Latitude        float32   `gorm:"not null"`
	RegionID        uuid.UUID `gorm:"not null"`

	Region    Region `gorm:"foreignKey:RegionID;references:ID"`
	CreatedBy User   `gorm:"foreignKey:CreatedByID;references:ID"`
}

type EventComment struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey"`
	EventID     uuid.UUID `gorm:"not null"`
	CreatedByID uuid.UUID `gorm:"not null"`
	CommentText string    `gorm:"not null"`

	LinkedEvent Event `gorm:"foreignKey:EventID;references:ID"`
	CreatedBy   User  `gorm:"foreignKey:CreatedByID;references:ID"`
}

type EventPhoto struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	EventID uuid.UUID `gorm:"not null"`
	Url     string    `gorm:"not null"`

	Event Event `gorm:"foreignKey:EventID;references:ID"`
}

func (ev *Event) BeforeDelete(tx *gorm.DB) error {
	var eventPhotos []EventPhoto
	// TODO: delete photos from dir.
	if err := tx.Table("event_photos").Where("event_id = ?", ev.ID).Find(&eventPhotos).Unscoped().Delete(&eventPhotos).Error; err != nil {
		return err
	}

	var eventComments []EventComment
	if err := tx.Table("event_comments").Where("event_id = ?", ev.ID).Find(&eventComments).Unscoped().Delete(&eventComments).Error; err != nil {
		return err
	}

	var linkedEvents []UserStoryLinkedEvent
	if err := tx.Table("user_story_linked_events").Where("event_id = ?", ev.ID).Find(&linkedEvents).Unscoped().Delete(&linkedEvents).Error; err != nil {
		return err
	}

	var likes []Liked
	if err := tx.Table("likeds").Where("event_id = ?", ev.ID).Find(&likes).Unscoped().Delete(&likes).Error; err != nil {
		return err
	}

	return nil
}

// ----- EventRepo methods -----

type EventRepoCrud interface {
	Create(dto dto.CreateEvent, u User, region Region) (Event, error)
	Get(id uuid.UUID) (Event, error)
	GetEvents(page int, count int) ([]Place, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, dto dto.UpdateEvent, region Region) error
}

type EventRepo struct {
	repo *PgRepo
}

func NewEventRepo(repo *PgRepo) *EventRepo {
	_ = repo.Database.AutoMigrate(&Event{}, &EventComment{}, &EventPhoto{})

	return &EventRepo{
		repo: repo,
	}
}

func (repo *EventRepo) Create(dto dto.CreateEvent, user User, region Region) (Event, error) {
	event := Event{
		ID:              uuid.New(),
		BeginDate:       time.Unix(dto.BeginDate, 0),
		EndDate:         time.Unix(dto.EndDate, 0),
		PaymentRequired: dto.PaymentNeed,
		CreatedByID:     user.ID,
		Title:           dto.Title,
		Description:     dto.Description,
		Longitude:       dto.Longitude,
		Latitude:        dto.Latitude,
		RegionID:        region.ID,

		CreatedBy: user,
		Region:    region,
	}

	return event, repo.repo.Database.Table("events").Create(&event).Error
}

func (repo *EventRepo) Get(id uuid.UUID) (event Event, err error) {
	return event, repo.repo.Database.Preload("CreatedBy").Preload("Region").Where("id = ?", id).First(&event).Error
}

func (repo *EventRepo) GetEvents(page int, count int) (events []Event, err error) {
	switch {
	case count > 15:
		count = 15
	case count <= 0:
		count = 15
	}
	offset := (page - 1) * count

	return events, repo.repo.Database.Preload("CreatedBy").Preload("Region").Offset(offset).Limit(count).Find(&events).Error
}

func (repo *EventRepo) Delete(id uuid.UUID) error {
	return repo.repo.Database.Unscoped().Delete(&Event{ID: id}).Error
}

func (repo *EventRepo) Update(id uuid.UUID, dto dto.UpdateEvent, region Region) error {
	event, err := repo.Get(id)
	if err != nil {
		return err
	}

	event.BeginDate = time.Unix(dto.BeginDate, 0)
	event.EndDate = time.Unix(dto.EndDate, 0)
	event.PaymentRequired = dto.PaymentNeed
	event.Title = dto.Title
	event.Description = dto.Description
	event.Longitude = dto.Longitude
	event.Latitude = dto.Latitude
	event.RegionID = region.ID
	event.Region = region

	return repo.repo.Database.Save(&event).Error
}

// ----- EventRepo methods: EventPhoto -----

func (repo *EventRepo) GetImages(id uuid.UUID) ([]string, error) {
	var photos []EventPhoto
	err := repo.repo.Database.Where("event_id = ?", id).Find(&photos).Error

	var result []string
	for _, photo := range photos {
		result = append(result, photo.Url)
	}

	return result, err
}

func (repo *EventRepo) CreateImages(id uuid.UUID, imgUrl string) error {
	eventPhoto := EventPhoto{
		ID:      uuid.New(),
		EventID: id,
		Url:     imgUrl,
	}

	return repo.repo.Database.Create(&eventPhoto).Error
}

func (repo *EventRepo) DeleteImages(url string) error {
	var photos []EventPhoto
	repo.repo.Database.Where("url = ?", url).Find(&photos)
	return repo.repo.Database.Table("event_photos").Unscoped().Delete(&photos).Error
}

// ----- EventRepo methods: EventComment -----

func (repo *EventRepo) CreateComment(commentDto dto.CreateEventCommentDto, user User, event Event) (EventComment, error) {
	eventComment := EventComment{
		ID:          uuid.New(),
		EventID:     event.ID,
		CreatedByID: user.ID,
		CommentText: commentDto.CommentBody,
		LinkedEvent: event,
		CreatedBy:   user,
	}

	return eventComment, repo.repo.Database.Create(&eventComment).Error
}

func (repo *EventRepo) GetComments(eventId uuid.UUID, page int, count int) (comments []EventComment, err error) {
	switch {
	case count > 15:
		count = 15
	case count <= 0:
		count = 15
	}
	offset := (page - 1) * count

	return comments, repo.repo.Database.Preload("CreatedBy").Preload("LinkedEvent").Offset(offset).Limit(count).
		Where("event_id = ?", eventId).Order("created_at DESC").Find(&comments).Error
}

func (repo *EventRepo) GetCommentByID(commentId uuid.UUID) (comment EventComment, err error) {
	return comment, repo.repo.Database.Preload("CreatedBy").Where("id = ?", commentId).Take(&comment).Error
}

func (repo *EventRepo) DeleteComments(commentId uuid.UUID) error {
	var comment EventComment
	repo.repo.Database.Where("id = ?", commentId).Find(&comment)
	return repo.repo.Database.Table("event_comments").Unscoped().Delete(&comment).Error
}

func (repo *EventRepo) UpdateComment(commentId uuid.UUID, updateDto dto.UpdateEventCommentDto) error {
	var comment EventComment
	repo.repo.Database.Where("id = ?", commentId).Find(&comment)
	comment.CommentText = updateDto.CommentBody
	return repo.repo.Database.Table("event_comments").Save(&comment).Error
}

// ----- Conversations -----

func EventToEvent(event Event, photoUrls []string) dto.EventDto {
	return dto.EventDto{
		ID:          event.ID,
		BeginDate:   event.BeginDate.Unix(),
		EndDate:     event.EndDate.Unix(),
		PaymentNeed: event.PaymentRequired,
		Title:       event.Title,
		Description: event.Description,
		Longitude:   event.Longitude,
		Latitude:    event.Latitude,
		CreatedBy:   UserToBaseUser(event.CreatedBy),
		RegionInfo:  RegionToRegion(event.Region),
		EventPhotos: photoUrls,
	}
}

func EventToUpdateEvent(event Event) dto.UpdateEvent {
	return dto.UpdateEvent{
		BeginDate:   event.BeginDate.Unix(),
		EndDate:     event.EndDate.Unix(),
		PaymentNeed: event.PaymentRequired,
		Title:       event.Title,
		Description: event.Description,
		Longitude:   event.Longitude,
		Latitude:    event.Latitude,
		RegionID:    event.Region.RegionShortName,
	}
}

func CommentToComment(comment EventComment) dto.EventCommentDto {
	return dto.EventCommentDto{
		ID:            comment.ID.String(),
		CreatedBy:     UserToBaseUser(comment.CreatedBy),
		LinkedEventID: comment.LinkedEvent.ID.String(),
		CommentBody:   comment.CommentText,
		UpdatedAt:     comment.UpdatedAt.Unix(),
		CreatedAt:     comment.CreatedAt.Unix(),
	}
}
