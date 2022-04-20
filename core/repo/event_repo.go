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
	RegionID        uuid.UUID

	region    Region `gorm:"foreignKey:region_id;references:id"`
	createdBy User   `gorm:"foreignKey:user_id;references:id"`
}

type EventComment struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey"`
	EventID     uuid.UUID `gorm:"not null"`
	CreatedByID uuid.UUID `gorm:"not null"`
	CommentText string    `gorm:"not null"`

	LinkedEvent Event `gorm:"foreignKey:event_id;references:id"`
	CreatedBy   User  `gorm:"foreignKey:created_by_id;references:id"`
}

type EventPhoto struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	EventID uuid.UUID `gorm:"not null"`
	Url     string    `gorm:"not null"`
	Size    int       `gorm:"not null"`

	event Event
}

// ----- EventRepo methods -----

type EventRepoCrud interface {
	Create(dto dto.CreateEvent, u User) (Event, error)
}

type EventRepo struct {
	Repo *PgRepo
}

func NewEventRepo(repo *PgRepo) *EventRepo {
	_ = repo.Database.AutoMigrate(&Event{}, &EventComment{}, &EventPhoto{})

	return &EventRepo{
		Repo: repo,
	}
}

func (repo *EventRepo) Create(dto dto.CreateEvent, u User) (Event, error) {
	event := Event{
		ID:              uuid.New(),
		BeginDate:       dto.BeginDate,
		EndDate:         dto.EndDate,
		PaymentRequired: false,
		CreatedByID:     u.ID,
		createdBy:       u,
		Title:           dto.Title,
		Description:     dto.Description,
		Longitude:       dto.Longitude,
		Latitude:        dto.Latitude,
	}

	return event, repo.Repo.Database.Create(&event).Error
}

//u, _ := userRepo.GetByEmail("zindoplay@gmail.com")
//e, _ := eventRepo.Create(dto.CreateEvent{
//	BeginDate:       time.Now(),
//	EndDate:         time.Now(),
//	PaymentRequired: false,
//	Title:           "",
//	Description:     "",
//	Longitude:       0,
//	Latitude:        0,
//}, u)
//eventRepo.Repo.Database.Create(&repo.EventComment{
//	ID:          uuid.New(),
//	EventID:     e.ID,
//	CreatedByID: u.ID,
//	CommentText: "qwerty123123",
//	LinkedEvent: e,
//	CreatedBy:   u,
//})
