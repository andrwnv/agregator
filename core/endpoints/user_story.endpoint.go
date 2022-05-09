package endpoints

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type UserStoryEndpoint struct {
	userStoryRepo *repo.UserStoryRepo
	userEndpoint  *UserEndpoint
	eventEndpoint *EventEndpoint
	placeEndpoint *PlaceEndpoint
}

func NewUserStoryEndpoint(userStoryRepo *repo.UserStoryRepo,
	userEndpoint *UserEndpoint,
	eventEndpoint *EventEndpoint,
	placeEndpoint *PlaceEndpoint) *UserStoryEndpoint {

	return &UserStoryEndpoint{
		userStoryRepo: userStoryRepo,
		userEndpoint:  userEndpoint,
		eventEndpoint: eventEndpoint,
		placeEndpoint: placeEndpoint,
	}
}

func (e *UserStoryEndpoint) GetFullStory(id uuid.UUID) (repo.UserStory, error) {
	return e.userStoryRepo.GetStoryByID(id)
}

func (e *UserStoryEndpoint) Get(id uuid.UUID) Result {
	story, err := e.userStoryRepo.GetStoryByID(id)
	events, _ := e.userStoryRepo.GetLinkedEvent(id)
	places, _ := e.userStoryRepo.GetLinkedPlace(id)
	photos, _ := e.userStoryRepo.GetImages(id)

	return Result{repo.StoryToStory(story, events, places, photos), err}
}

func (e *UserStoryEndpoint) Create(
	createDto dto.CreateUserStoryDto,
	eventsId []uuid.UUID,
	placesId []uuid.UUID,
	userInfo dto.BaseUserInfo) Result {

	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find user for create story.")}
	}

	var events []repo.Event
	for _, id := range eventsId {
		event, extractErr := e.eventEndpoint.GetFullEvent(id)
		if extractErr == nil {
			events = append(events, event)
		}
	}

	var places []repo.Place
	for _, id := range placesId {
		place, extractErr := e.placeEndpoint.GetFullPlace(id)
		if extractErr == nil {
			places = append(places, place)
		}
	}

	story, err := e.userStoryRepo.Create(createDto, user, events, places)
	linkedEvents, _ := e.userStoryRepo.GetLinkedEvent(story.ID)
	linkedPlaces, _ := e.userStoryRepo.GetLinkedPlace(story.ID)
	return Result{repo.StoryToStory(story, linkedEvents, linkedPlaces, []string{}), err}
}

func (e *UserStoryEndpoint) Delete(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your story!")}
	}

	err = e.userStoryRepo.Delete(id)
	return Result{err != nil, err}
}

func (e *UserStoryEndpoint) Update(id uuid.UUID, updateDto dto.UpdateUserStoryDto, userInfo dto.BaseUserInfo) Result {
	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your story!")}
	}

	err = e.userStoryRepo.Update(id, updateDto)
	return Result{err != nil, err}
}

// ----- UserStoryEndpoint: Linked images -----

func (e *UserStoryEndpoint) UpdateImages(id uuid.UUID, userInfo dto.BaseUserInfo,
	filesToCreate []string, filesToDelete []string) Result {

	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your story!")}
	}

	for _, url := range filesToCreate {
		err := e.userStoryRepo.CreateImages(story.ID, url)
		if err != nil {
			return Result{false, err}
		}
	}

	for _, url := range filesToDelete {
		err := e.userStoryRepo.DeleteImages(url)
		// TODO: delete photos from dir.
		if err != nil {
			return Result{false, err}
		}
	}

	return Result{true, nil}
}

// ----- UserStoryEndpoint: Linked events -----

func (e *UserStoryEndpoint) UpdateLinkedEvents(id uuid.UUID, userInfo dto.BaseUserInfo,
	toCreate []string, toDelete []string) Result {

	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your story!")}
	}

	for _, eventId := range toCreate {
		fullEvent, getErr := e.eventEndpoint.GetFullEvent(uuid.MustParse(eventId))
		if getErr == nil {
			err := e.userStoryRepo.AddLinkedEvent(story.ID, fullEvent)
			if err != nil {
				return Result{false, err}
			}
		}
	}

	for _, eventId := range toDelete {
		err := e.userStoryRepo.DeleteLinkedEvent(story.ID, uuid.MustParse(eventId))
		if err != nil {
			return Result{false, err}
		}
	}

	return Result{true, nil}
}

// ----- UserStoryEndpoint: Linked places -----

func (e *UserStoryEndpoint) UpdateLinkedPlaces(id uuid.UUID, userInfo dto.BaseUserInfo,
	toCreate []string, toDelete []string) Result {

	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your story!")}
	}

	for _, placeId := range toCreate {
		fullPlace, getErr := e.placeEndpoint.GetFullPlace(uuid.MustParse(placeId))
		if getErr == nil {
			err := e.userStoryRepo.AddLinkedPlace(story.ID, fullPlace)
			if err != nil {
				return Result{false, err}
			}
		}
	}

	for _, placeId := range toDelete {
		err := e.userStoryRepo.DeleteLinkedPlace(story.ID, uuid.MustParse(placeId))
		if err != nil {
			return Result{false, err}
		}
	}

	return Result{true, nil}
}
