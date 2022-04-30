package repo

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Region struct {
	gorm.Model

	ID              uuid.UUID `gorm:"primaryKey"`
	RegionName      string    `gorm:"not null"`
	RegionShortName string    `gorm:"not null"`
}

// ----- RegionRepo methods -----

type RegionRepo struct {
	repo    *PgRepo
	regions map[string]string
}

func NewRegionRepo(repo *PgRepo) *RegionRepo {
	_ = repo.Database.AutoMigrate(&Region{})

	rep := &RegionRepo{
		repo: repo,
		regions: map[string]string{
			"EG": "Egypt",
			"IN": "India",
			"JP": "Japan",
			"RU": "Russian Federation",
		},
	}
	rep.initBaseRegions()

	return rep
}

func (r *RegionRepo) GetByRegionID(regionID string) (region Region, err error) {
	err = r.repo.Database.Where(&Region{RegionShortName: regionID}).First(&region).Error
	return region, err
}

func (r *RegionRepo) GetByRegionName(regionName string) (region Region, err error) {
	err = r.repo.Database.Where(&Region{RegionName: regionName}).First(&region).Error
	return region, err
}

func (r *RegionRepo) initBaseRegions() {
	for key, value := range r.regions {
		err := r.repo.Database.Where(&Region{RegionShortName: key, RegionName: value}).First(&Region{}).Error
		if err != nil {
			createErr := r.repo.Database.Create(&Region{
				ID:              uuid.New(),
				RegionName:      value,
				RegionShortName: key,
			}).Error

			if createErr != nil {
				misc.ReportCritical("Cant get info from db")
				return
			}
		}
	}
}

// ----- Conversations -----

func RegionToRegion(region Region) dto.RegionDto {
	return dto.RegionDto{
		RegionID:   region.RegionShortName,
		RegionName: region.RegionName,
	}
}
