package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	gorm.Model

	ID         uuid.UUID `gorm:"primaryKey"`
	RegionName string    `gorm:"not null"`
	regionID   string    `gorm:"not null"`
}

// ----- RegionRepo methods -----

type RegionRepo struct {
	Repo *PgRepo
}

func NewRegionRepo(repo *PgRepo) *RegionRepo {
	_ = repo.Database.AutoMigrate(&Region{})

	return &RegionRepo{
		Repo: repo,
	}
}
