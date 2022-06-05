package repo

import (
	"errors"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Liked struct {
	gorm.Model

	ID      uuid.UUID  `gorm:"primaryKey"`
	UserID  uuid.UUID  `gorm:"not null"`
	EventID *uuid.UUID `gorm:"null"`
	PlaceID *uuid.UUID `gorm:"null"`

	CreatedBy   User   `gorm:"foreignKey:UserID;references:ID"`
	LinkedEvent *Event `gorm:"foreignKey:EventID;references:ID"`
	LinkedPlace *Place `gorm:"foreignKey:PlaceID;references:ID"`
}

// ----- LikedRepo methods -----

type LikedRepo struct {
	repo *PgRepo
}

func NewLikedRepo(repo *PgRepo) *LikedRepo {
	_ = repo.Database.AutoMigrate(&Liked{})

	rep := &LikedRepo{
		repo: repo,
	}

	return rep
}

func (repo *LikedRepo) Get(user dto.BaseUserInfo, page int, count int) (likedList []Liked, err error) {
	switch {
	case count > 15:
		count = 15
	case count <= 0:
		count = 15
	}
	offset := (page - 1) * count

	return likedList, repo.repo.Database.Preload("CreatedBy").Preload("LinkedEvent").Preload("LinkedPlace").
		Offset(offset).Limit(count).Where("user_id = ?", user.ID).Find(&likedList).Error
}

func (repo *LikedRepo) Dislike(user User, id uuid.UUID) error {
	like := Liked{}
	repo.repo.Database.Where("event_id = ?", id).Or("place_id = ?", id).Take(&like)
	return repo.repo.Database.Unscoped().Delete(&like).Error
}

func (repo *LikedRepo) LikeEvent(user User, event Event) (Liked, error) {
	var likedList []Liked
	repo.repo.Database.Where("user_id = ?", user.ID).Where("event_id = ?", event.ID).Find(&likedList)
	if len(likedList) > 0 {
		return Liked{}, errors.New("already liked")
	}

	liked := Liked{
		ID:          uuid.New(),
		UserID:      uuid.UUID{},
		EventID:     &event.ID,
		PlaceID:     nil,
		CreatedBy:   user,
		LinkedEvent: &event,
		LinkedPlace: nil,
	}
	return liked, repo.repo.Database.Create(&liked).Error
}

func (repo *LikedRepo) LikePlace(user User, place Place) (Liked, error) {
	var likedList []Liked
	repo.repo.Database.Where("user_id = ?", user.ID).Where("place_id = ?", place.ID).Find(&likedList)
	if len(likedList) > 0 {
		return Liked{}, errors.New("already liked")
	}

	liked := Liked{
		ID:          uuid.New(),
		UserID:      uuid.UUID{},
		EventID:     nil,
		PlaceID:     &place.ID,
		CreatedBy:   user,
		LinkedEvent: nil,
		LinkedPlace: &place,
	}
	return liked, repo.repo.Database.Create(&liked).Error
}

// ----- Conversations -----

func LikeToLike(liked Liked) dto.LikedDto {
	return dto.LikedDto{
		User:    UserToBaseUser(liked.CreatedBy),
		EventID: liked.EventID,
		PlaceID: liked.PlaceID,
	}
}
