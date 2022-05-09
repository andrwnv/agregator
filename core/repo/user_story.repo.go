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
	CreatedBy    User      `gorm:"foreignKey:CreatedByID;references:ID"`
}

type UserStoryLinkedEvent struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey"`
	EventID     uuid.UUID `gorm:"not null"`
	StoryID     uuid.UUID `gorm:"not null"`
	Story       UserStory `gorm:"foreignKey:StoryID;references:ID"`
	LinkedEvent Event     `gorm:"foreignKey:EventID;references:ID"`
}

type UserStoryLinkedPlace struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey"`
	PlaceID     uuid.UUID `gorm:"not null"`
	StoryID     uuid.UUID `gorm:"not null"`
	Story       UserStory `gorm:"foreignKey:StoryID;references:ID"`
	LinkedPlace Place     `gorm:"foreignKey:PlaceID;references:ID"`
}

type UserStoryLinkedPhoto struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	StoryID uuid.UUID `gorm:"not null"`
	Url     string    `gorm:"not null"`
	Story   UserStory `gorm:"foreignKey:StoryID;references:ID"`
}

func (us *UserStory) BeforeDelete(tx *gorm.DB) error {
	var linkedEvents []UserStoryLinkedEvent
	if err := tx.Table("user_story_linked_events").Where("story_id = ?", us.ID).Find(&linkedEvents).Unscoped().Delete(&linkedEvents).Error; err != nil {
		return err
	}

	var linkedPlaces []UserStoryLinkedPlace
	if err := tx.Table("user_story_linked_places").Where("story_id = ?", us.ID).Find(&linkedPlaces).Unscoped().Delete(&linkedPlaces).Error; err != nil {
		return err
	}

	var linkedPhotos []UserStoryLinkedPhoto
	if err := tx.Table("user_story_linked_photos").Where("story_id = ?", us.ID).Find(&linkedPhotos).Unscoped().Delete(&linkedPhotos).Error; err != nil {
		return err
	}

	return nil
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
