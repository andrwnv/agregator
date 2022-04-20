package repo

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Place struct {
	gorm.Model

	ID              uuid.UUID `gorm:"primaryKey"`
	PaymentRequired bool      `gorm:"default: false"`
	CreatedByID     uuid.UUID `gorm:"not null"`
	Title           string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	Longitude       float32   `gorm:"not null"`
	Latitude        float32   `gorm:"not null"`

	createdBy User
}

type PlaceComment struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey"`
	EventID     uuid.UUID `gorm:"not null"`
	CreatedByID uuid.UUID `gorm:"not null"`
	CommentText string    `gorm:"not null"`

	linkedEvent Place `gorm:"foreignKey:event_id;references:id"`
	createdBy   User  `gorm:"foreignKey:created_by_id;references:id"`
}

type PlacePhoto struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	EventID uuid.UUID `gorm:"not null"`
	Url     string    `gorm:"not null"`
	Size    int       `gorm:"not null"`

	event Event
}

// ----- PlaceRepo methods -----

type PlaceRepoCrud interface {
	Create(dto dto.CreatePlace, u User) (Place, error)
}

type PlaceRepo struct {
	Repo *PgRepo
}

func NewPlaceRepo(repo *PgRepo) *PlaceRepo {
	_ = repo.Database.AutoMigrate(&Place{}, &PlaceComment{}, &PlacePhoto{})

	return &PlaceRepo{
		Repo: repo,
	}
}

func (repo *PlaceRepo) Create(dto dto.CreatePlace, u User) (Place, error) {
	place := Place{
		ID:              uuid.New(),
		PaymentRequired: false,
		CreatedByID:     u.ID,
		createdBy:       u,
		Title:           dto.Title,
		Description:     dto.Description,
		Longitude:       dto.Longitude,
		Latitude:        dto.Latitude,
	}

	return place, repo.Repo.Database.Create(&place).Error
}
