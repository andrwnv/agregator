package repo

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model

	ID              uuid.UUID  `gorm:"primaryKey"`
	BeginDate       time.Time  `gorm:"not null"`
	EndDate         time.Time  `gorm:"not null"`
	PaymentRequired bool       `gorm:"default: false"`
	CreatedByID     uuid.UUID  `gorm:"not null"`
	Title           string     `gorm:"not null"`
	Description     string     `gorm:"not null"`
	Longitude       float32    `gorm:"not null"`
	Latitude        float32    `gorm:"not null"`
	RegionID        *uuid.UUID `gorm:"null"`

	Region    *Region `gorm:"foreignKey:RegionID;references:ID"`
	CreatedBy User    `gorm:"foreignKey:CreatedByID;references:ID"`
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
	Size    int       `gorm:"not null"`

	Event Event `gorm:"foreignKey:EventID;references:ID"`
}

// ----- EventRepo methods -----

type EventRepoCrud interface {
	Create(dto dto.CreateEvent, u User) (Event, error)
	Get(id uuid.UUID) (Event, error)
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

func (repo *EventRepo) Create(dto dto.CreateEvent, u User) (Event, error) {
	event := Event{
		ID:              uuid.New(),
		BeginDate:       time.Unix(dto.BeginDate, 0),
		EndDate:         time.Unix(dto.EndDate, 0),
		PaymentRequired: false,
		CreatedByID:     u.ID,
		CreatedBy:       u,
		Title:           dto.Title,
		Description:     dto.Description,
		Longitude:       dto.Longitude,
		Latitude:        dto.Latitude,
	}

	return event, repo.repo.Database.Create(&event).Error
}

func (repo *EventRepo) Get(id uuid.UUID) (event Event, err error) {
	err = repo.repo.Database.Preload("CreatedBy").Where("id = ?", id).First(&event).Error
	return event, err
}

// ----- Conversations -----

func EventToEvent(event Event) dto.EventDto {
	return dto.EventDto{
		ID:              event.ID,
		BeginDate:       event.BeginDate.Unix(),
		EndDate:         event.EndDate.Unix(),
		PaymentRequired: event.PaymentRequired,
		Title:           event.Title,
		Description:     event.Description,
		Longitude:       event.Longitude,
		Latitude:        event.Latitude,
		CreatedBy:       UserToBaseUser(event.CreatedBy),
	}
}
