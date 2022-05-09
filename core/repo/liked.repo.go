package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Liked struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	UserID  uuid.UUID `gorm:"not null"`
	EventID uuid.UUID `gorm:"null"`
	PlaceID uuid.UUID `gorm:"null"`

	CreatedBy   User
	LinkedEvent Event
	LinkedPlace Place
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
