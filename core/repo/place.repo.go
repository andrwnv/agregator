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
	RegionID        uuid.UUID

	Region    Region `gorm:"foreignKey:RegionID;references:ID"`
	CreatedBy User   `gorm:"foreignKey:CreatedByID;references:ID"`
}

type PlaceComment struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey"`
	PlaceID     uuid.UUID `gorm:"not null"`
	CreatedByID uuid.UUID `gorm:"not null"`
	CommentText string    `gorm:"not null"`

	LinkedPlace Place `gorm:"foreignKey:PlaceID;references:ID"`
	CreatedBy   User  `gorm:"foreignKey:CreatedByID;references:ID"`
}

type PlacePhoto struct {
	gorm.Model

	ID      uuid.UUID `gorm:"primaryKey"`
	PlaceID uuid.UUID `gorm:"not null"`
	Url     string    `gorm:"not null"`
	Size    int       `gorm:"not null"`

	LinkedPlace Place `gorm:"foreignKey:PlaceID;references:ID"`
}

func (ev *Place) BeforeDelete(tx *gorm.DB) error {
	var placePhotos []PlacePhoto
	// TODO: delete photos from dir.
	if err := tx.Table("place_photos").Where("place_id = ?", ev.ID).Find(&placePhotos).Unscoped().Delete(&placePhotos).Error; err != nil {
		return err
	}

	var placeComments []PlaceComment
	if err := tx.Table("place_comments").Where("place_id = ?", ev.ID).Find(&placeComments).Unscoped().Delete(&placeComments).Error; err != nil {
		return err
	}

	var linkedPlaces []UserStoryLinkedPlace
	if err := tx.Table("user_story_linked_places").Where("place_id = ?", ev.ID).Find(&linkedPlaces).Unscoped().Delete(&linkedPlaces).Error; err != nil {
		return err
	}

	var likes []Liked
	if err := tx.Table("likeds").Where("place_id = ?", ev.ID).Find(&likes).Unscoped().Delete(&likes).Error; err != nil {
		return err
	}

	return nil
}

// ----- PlaceRepo methods -----

type PlaceRepoCrud interface {
	Create(dto dto.CreatePlace, u User, region Region) (Place, error)
	Get(id uuid.UUID) (Place, error)
	GetPlaces(page int, count int) ([]Place, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, dto dto.UpdatePlace, region Region) error
}

type PlaceRepo struct {
	repo *PgRepo
}

func NewPlaceRepo(repo *PgRepo) *PlaceRepo {
	_ = repo.Database.AutoMigrate(&Place{}, &PlaceComment{}, &PlacePhoto{})

	return &PlaceRepo{
		repo: repo,
	}
}

func (repo *PlaceRepo) GetPlaces(page int, count int) (places []Place, err error) {
	switch {
	case count > 15:
		count = 15
	case count <= 0:
		count = 15
	}
	offset := (page - 1) * count

	return places, repo.repo.Database.Preload("CreatedBy").Preload("Region").Offset(offset).Limit(count).Find(&places).Error
}

func (repo *PlaceRepo) Create(dto dto.CreatePlace, u User, region Region) (Place, error) {
	place := Place{
		ID:              uuid.New(),
		PaymentRequired: false,
		CreatedByID:     u.ID,
		CreatedBy:       u,
		Title:           dto.Title,
		Description:     dto.Description,
		Longitude:       dto.Longitude,
		Latitude:        dto.Latitude,
		RegionID:        region.ID,
		Region:          region,
	}

	return place, repo.repo.Database.Create(&place).Error
}

func (repo *PlaceRepo) Get(id uuid.UUID) (place Place, err error) {
	return place, repo.repo.Database.Preload("CreatedBy").Preload("Region").Where("id = ?", id).First(&place).Error
}

func (repo *PlaceRepo) Delete(id uuid.UUID) error {
	return repo.repo.Database.Unscoped().Delete(&Place{ID: id}).Error
}

func (repo *PlaceRepo) Update(id uuid.UUID, dto dto.UpdatePlace, region Region) error {
	place, err := repo.Get(id)
	if err != nil {
		return err
	}

	place.PaymentRequired = dto.PaymentNeed
	place.Title = dto.Title
	place.Description = dto.Description
	place.Longitude = dto.Longitude
	place.Latitude = dto.Latitude
	place.RegionID = region.ID
	place.Region = region

	return repo.repo.Database.Save(&place).Error
}

// ----- PlaceRepo methods: PlacePhoto -----

func (repo *PlaceRepo) GetImages(id uuid.UUID) ([]string, error) {
	var photos []PlacePhoto
	err := repo.repo.Database.Where("place_id = ?", id).Find(&photos).Error

	var result []string
	for _, photo := range photos {
		result = append(result, photo.Url)
	}

	return result, err
}

func (repo *PlaceRepo) CreateImages(id uuid.UUID, imgUrl string) error {
	placePhoto := PlacePhoto{
		ID:      uuid.New(),
		PlaceID: id,
		Url:     imgUrl,
	}

	return repo.repo.Database.Create(&placePhoto).Error
}

func (repo *PlaceRepo) DeleteImages(url string) error {
	var photos []PlacePhoto
	repo.repo.Database.Where("url = ?", url).Find(&photos)
	return repo.repo.Database.Table("place_photos").Unscoped().Delete(&photos).Error
}

// ----- PlaceRepo methods: PlaceComment -----

func (repo *PlaceRepo) CreateComment(commentDto dto.CreatePlaceCommentDto, user User, place Place) (PlaceComment, error) {
	placeComment := PlaceComment{
		ID:          uuid.New(),
		PlaceID:     place.ID,
		CreatedByID: user.ID,
		CommentText: commentDto.CommentBody,
		LinkedPlace: place,
		CreatedBy:   user,
	}

	return placeComment, repo.repo.Database.Create(&placeComment).Error
}

func (repo *PlaceRepo) GetComments(placeId uuid.UUID, page int, count int) (comments []PlaceComment, err error) {
	switch {
	case count > 15:
		count = 15
	case count <= 0:
		count = 15
	}
	offset := (page - 1) * count

	return comments, repo.repo.Database.Preload("CreatedBy").Preload("LinkedPlace").Offset(offset).Limit(count).
		Where("place_id = ?", placeId).Find(&comments).Error
}

func (repo *PlaceRepo) GetCommentByID(commentId uuid.UUID) (comment PlaceComment, err error) {
	return comment, repo.repo.Database.Preload("CreatedBy").Where("id = ?", commentId).Take(&comment).Error
}

func (repo *PlaceRepo) DeleteComments(commentId uuid.UUID) error {
	var comment PlaceComment
	repo.repo.Database.Where("id = ?", commentId).Find(&comment)
	return repo.repo.Database.Table("place_comments").Unscoped().Delete(&comment).Error
}

func (repo *PlaceRepo) UpdateComment(commentId uuid.UUID, updateDto dto.UpdatePlaceCommentDto) error {
	var comment PlaceComment
	repo.repo.Database.Where("id = ?", commentId).Find(&comment)
	comment.CommentText = updateDto.CommentBody
	return repo.repo.Database.Table("place_comments").Save(&comment).Error
}

// ----- Conversations -----

func PlaceToPlace(place Place, photoUrls []string) dto.PlaceDto {
	return dto.PlaceDto{
		ID:          place.ID,
		PaymentNeed: place.PaymentRequired,
		Title:       place.Title,
		Description: place.Description,
		Longitude:   place.Longitude,
		Latitude:    place.Latitude,
		CreatedBy:   UserToBaseUser(place.CreatedBy),
		RegionInfo:  RegionToRegion(place.Region),
		PlacePhotos: photoUrls,
	}
}

func PlaceToUpdatePlace(place Place) dto.UpdatePlace {
	return dto.UpdatePlace{
		PaymentNeed: place.PaymentRequired,
		Title:       place.Title,
		Description: place.Description,
		Longitude:   place.Longitude,
		Latitude:    place.Latitude,
		RegionID:    place.Region.RegionShortName,
	}
}

func CommentToCommentDto(comment PlaceComment) dto.PlaceCommentDto {
	return dto.PlaceCommentDto{
		ID:            comment.ID.String(),
		CreatedBy:     UserToBaseUser(comment.CreatedBy),
		LinkedPlaceID: comment.LinkedPlace.ID.String(),
		CommentBody:   comment.CommentText,
		UpdatedAt:     comment.UpdatedAt.Unix(),
		CreatedAt:     comment.CreatedAt.Unix(),
	}
}
