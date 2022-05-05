package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStory struct {
	gorm.Model

	ID           uuid.UUID `gorm:"primaryKey"`
	CreatedByID  uuid.UUID `gorm:"not null"`
	Title        string    `gorm:"not null"`
	LongReadText string    `gorm:"not null"`

	createdBy User
}

type UserStoryLinkedEvent struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	EventID uuid.UUID `gorm:"not null"`
	StoryID uuid.UUID `gorm:"not null"`

	linkedEvent Place `gorm:"foreignKey:event_id;references:id"`
	createdBy   User  `gorm:"foreignKey:created_by_id;references:id"`
}

type UserStoryLinkedPlace struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	PlaceID uuid.UUID `gorm:"not null"`
	StoryID uuid.UUID `gorm:"not null"`

	linkedEvent Place `gorm:"foreignKey:place_id;references:id"`
	createdBy   User  `gorm:"foreignKey:created_by_id;references:id"`
}

type UserStoryLinkedPhoto struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	StoryID uuid.UUID `gorm:"not null"`
	Url     string    `gorm:"not null"`
	Size    int       `gorm:"not null"`

	story UserStory
}

// ----- UserStoryRepo methods -----

type UserStoryRepo struct {
	Repo *PgRepo
}

func NewUserStoryRepo(repo *PgRepo) *UserStoryRepo {
	_ = repo.Database.AutoMigrate(&UserStory{},
		&UserStoryLinkedEvent{},
		&UserStoryLinkedPhoto{},
		&UserStoryLinkedPlace{})

	return &UserStoryRepo{
		Repo: repo,
	}
}
