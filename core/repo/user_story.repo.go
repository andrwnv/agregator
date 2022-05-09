package repo

import (
	"github.com/andrwnv/event-aggregator/core/dto"
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
	repo *PgRepo
}

func NewUserStoryRepo(repo *PgRepo) *UserStoryRepo {
	_ = repo.Database.AutoMigrate(&UserStory{},
		&UserStoryLinkedEvent{},
		&UserStoryLinkedPhoto{},
		&UserStoryLinkedPlace{})

	return &UserStoryRepo{
		repo: repo,
	}
}

// ----- UserStoryRepo methods: UserStory -----

func (repo *UserStoryRepo) Create(createDto dto.CreateUserStoryDto, user User, events []Event, places []Place) (UserStory, error) {
	userStory := UserStory{
		ID:           uuid.New(),
		CreatedByID:  user.ID,
		Title:        createDto.Title,
		LongReadText: createDto.LongReadText,
		CreatedBy:    user,
	}

	err := repo.repo.Database.Table("user_stories").Create(&userStory).Error
	if err != nil {
		return UserStory{}, err
	}

	if len(events) > 0 {
		var linkedEvents []UserStoryLinkedEvent
		for _, event := range events {
			linkedEvents = append(linkedEvents, UserStoryLinkedEvent{
				ID:          uuid.New(),
				EventID:     event.ID,
				StoryID:     userStory.ID,
				Story:       userStory,
				LinkedEvent: event,
			})
		}

		linkErr := repo.repo.Database.Table("user_story_linked_events").Create(&linkedEvents).Error
		if linkErr != nil {
			return userStory, linkErr
		}
	}

	if len(places) > 0 {
		var linkedPlaces []UserStoryLinkedPlace
		for _, place := range places {
			linkedPlaces = append(linkedPlaces, UserStoryLinkedPlace{
				ID:          uuid.New(),
				PlaceID:     place.ID,
				StoryID:     userStory.ID,
				Story:       userStory,
				LinkedPlace: place,
			})
		}

		linkErr := repo.repo.Database.Table("user_story_linked_places").Create(&linkedPlaces).Error
		if linkErr != nil {
			return userStory, linkErr
		}
	}

	return userStory, nil
}

func (repo *UserStoryRepo) Delete(id uuid.UUID) error {
	return repo.repo.Database.Unscoped().Delete(&UserStory{ID: id}).Error
}

func (repo *UserStoryRepo) GetStoryByID(id uuid.UUID) (story UserStory, err error) {
	return story, repo.repo.Database.Preload("CreatedBy").Where("id = ?", id).Find(&story).Error
}

func (repo *UserStoryRepo) GetStories(page int, count int) (stories []UserStory, err error) {
	switch {
	case count > 15:
		count = 15
	case count <= 0:
		count = 15
	}
	offset := (page - 1) * count

	return stories, repo.repo.Database.Preload("CreatedBy").Offset(offset).Limit(count).Find(&stories).Error
}

func (repo *UserStoryRepo) Update(id uuid.UUID, updateDto dto.UpdateUserStoryDto) error {
	story, err := repo.GetStoryByID(id)
	if err != nil {
		return err
	}

	story.Title = updateDto.Title
	story.LongReadText = updateDto.LongReadText

	return repo.repo.Database.Save(&story).Error
}

// ----- UserStoryRepo methods: LinkedEvent -----

func (repo *UserStoryRepo) AddLinkedEvent(id uuid.UUID, event Event) error {
	story, err := repo.GetStoryByID(id)
	if err != nil {
		return err
	}

	return repo.repo.Database.Create(&UserStoryLinkedEvent{
		ID:          uuid.New(),
		EventID:     event.ID,
		StoryID:     story.ID,
		Story:       story,
		LinkedEvent: event,
	}).Error
}

func (repo *UserStoryRepo) DeleteLinkedEvent(id uuid.UUID, eventId uuid.UUID) error {
	return repo.repo.Database.Unscoped().Table("user_story_linked_events").Where("event_id = ?", eventId).Where("story_id = ?", id).
		Delete(&UserStoryLinkedEvent{StoryID: id, EventID: eventId}).Error
}

func (repo *UserStoryRepo) GetLinkedEvent(id uuid.UUID) (result []UserStoryLinkedEvent, err error) {
	return result, repo.repo.Database.Preload("LinkedEvent").Preload("Story").
		Table("user_story_linked_events").Where("story_id = ?", id).Find(&result).Error
}

// ----- UserStoryRepo methods: LinkedPlace -----

func (repo *UserStoryRepo) AddLinkedPlace(id uuid.UUID, place Place) error {
	story, err := repo.GetStoryByID(id)
	if err != nil {
		return err
	}

	return repo.repo.Database.Create(&UserStoryLinkedPlace{
		ID:          uuid.New(),
		PlaceID:     place.ID,
		StoryID:     story.ID,
		Story:       story,
		LinkedPlace: place,
	}).Error
}

func (repo *UserStoryRepo) DeleteLinkedPlace(id uuid.UUID, placeId uuid.UUID) error {
	return repo.repo.Database.Unscoped().Table("user_story_linked_places").Where("place_id = ?", placeId).Where("story_id = ?", id).
		Delete(&UserStoryLinkedPlace{StoryID: id, PlaceID: placeId}).Error
}

func (repo *UserStoryRepo) GetLinkedPlace(id uuid.UUID) (result []UserStoryLinkedPlace, err error) {
	return result, repo.repo.Database.Preload("LinkedPlace").Preload("Story").
		Table("user_story_linked_places").Where("story_id = ?", id).Find(&result).Error
}

// ----- UserStoryRepo methods: Photo -----

func (repo *UserStoryRepo) GetImages(id uuid.UUID) ([]string, error) {
	var photos []UserStoryLinkedPhoto
	err := repo.repo.Database.Table("user_story_linked_photos").Where("story_id = ?", id).Find(&photos).Error

	var result []string
	for _, photo := range photos {
		result = append(result, photo.Url)
	}

	return result, err
}

func (repo *UserStoryRepo) CreateImages(id uuid.UUID, imgUrl string) error {
	storyPhoto := UserStoryLinkedPhoto{
		ID:      uuid.New(),
		StoryID: id,
		Url:     imgUrl,
	}

	return repo.repo.Database.Create(&storyPhoto).Error
}

func (repo *UserStoryRepo) DeleteImages(url string) error {
	var photos []UserStoryLinkedPhoto
	repo.repo.Database.Where("url = ?", url).Find(&photos)
	return repo.repo.Database.Table("user_story_linked_photos").Unscoped().Delete(&photos).Error
}

// ----- Conversations -----

func LinkedEventToLinkedEvent(linked UserStoryLinkedEvent) dto.LinkedEventDto {
	return dto.LinkedEventDto{
		ID:      linked.ID.String(),
		EventID: linked.LinkedEvent.ID.String(),
		StoryID: linked.StoryID.String(),
	}
}

func LinkedPlaceToLinkedPlace(linked UserStoryLinkedPlace) dto.LinkedPlaceDto {
	return dto.LinkedPlaceDto{
		ID:      linked.ID.String(),
		PlaceID: linked.LinkedPlace.ID.String(),
		StoryID: linked.StoryID.String(),
	}
}

func StoryToStory(
	story UserStory,
	events []UserStoryLinkedEvent,
	places []UserStoryLinkedPlace,
	photos []string) dto.UserStoryDto {

	var eventsDto []dto.LinkedEventDto
	var placesDto []dto.LinkedPlaceDto

	for _, event := range events {
		eventsDto = append(eventsDto, LinkedEventToLinkedEvent(event))
	}
	for _, place := range places {
		placesDto = append(placesDto, LinkedPlaceToLinkedPlace(place))
	}

	return dto.UserStoryDto{
		ID:           story.ID,
		Title:        story.Title,
		LongReadText: story.LongReadText,
		CreatedBy:    UserToBaseUser(story.CreatedBy),
		LinkedEvents: eventsDto,
		LinkedPlaces: placesDto,
		LinkedPhotos: photos,
	}
}
func StoryToUpdateStory(story UserStory) dto.UpdateUserStoryDto {
	return dto.UpdateUserStoryDto{
		Title:        story.Title,
		LongReadText: story.LongReadText,
	}
}
