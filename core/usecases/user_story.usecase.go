package usecases

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type UserStoryUsecase struct {
	userStoryRepo *repo.UserStoryRepo
	userUsecase   *UserUsecase
	eventUsecase  *EventUsecase
	placeUsecase  *PlaceUsecase
}

func NewUserStoryUsecase(userStoryRepo *repo.UserStoryRepo,
	userUsecase *UserUsecase,
	eventUsecase *EventUsecase,
	placeUsecase *PlaceUsecase) *UserStoryUsecase {

	return &UserStoryUsecase{
		userStoryRepo: userStoryRepo,
		userUsecase:   userUsecase,
		eventUsecase:  eventUsecase,
		placeUsecase:  placeUsecase,
	}
}

func (e *UserStoryUsecase) GetFullStory(id uuid.UUID) (repo.UserStory, error) {
	return e.userStoryRepo.GetStoryByID(id)
}

func (e *UserStoryUsecase) Get(id uuid.UUID) Result {
	story, err := e.userStoryRepo.GetStoryByID(id)
	events, _ := e.userStoryRepo.GetLinkedEvent(id)
	places, _ := e.userStoryRepo.GetLinkedPlace(id)
	photos, _ := e.userStoryRepo.GetImages(id)

	return Result{repo.StoryToStory(story, events, places, photos), err}
}

func (e *UserStoryUsecase) GetPaginated(page int, count int) Result {
	stories, err := e.userStoryRepo.GetStories(page, count)
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant extract stories.")}
	}

	total, err := e.userStoryRepo.GetTotalCount()
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant extract total count of stories.")}
	}

	var result []dto.ShortStoryInfoDto
	for _, story := range stories {
		result = append(result, dto.ShortStoryInfoDto{
			ID:           story.ID,
			Title:        story.Title,
			LongReadText: story.LongReadText,
		})
	}

	return Result{
		Value: dto.ShortStoryListDto{
			Page:      int64(page),
			ListSize:  int64(count),
			TotalSize: total,
			List:      result,
		},
		Error: nil,
	}
}

func (e *UserStoryUsecase) Create(
	createDto dto.CreateUserStoryDto,
	eventsId []uuid.UUID,
	placesId []uuid.UUID,
	userInfo dto.BaseUserInfo) Result {

	user, err := e.userUsecase.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeUsecaseError("Cant find user for create story.")}
	}

	var events []repo.Event
	for _, id := range eventsId {
		event, extractErr := e.eventUsecase.GetFullEvent(id)
		if extractErr == nil {
			events = append(events, event)
		}
	}

	var places []repo.Place
	for _, id := range placesId {
		place, extractErr := e.placeUsecase.GetFullPlace(id)
		if extractErr == nil {
			places = append(places, place)
		}
	}

	story, err := e.userStoryRepo.Create(createDto, user, events, places)
	linkedEvents, _ := e.userStoryRepo.GetLinkedEvent(story.ID)
	linkedPlaces, _ := e.userStoryRepo.GetLinkedPlace(story.ID)
	return Result{repo.StoryToStory(story, linkedEvents, linkedPlaces, []string{}), err}
}

func (e *UserStoryUsecase) Delete(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your story!")}
	}

	err = e.userStoryRepo.Delete(id)
	return Result{err != nil, err}
}

func (e *UserStoryUsecase) Update(id uuid.UUID, updateDto dto.UpdateUserStoryDto, userInfo dto.BaseUserInfo) Result {
	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your story!")}
	}

	err = e.userStoryRepo.Update(id, updateDto)
	return Result{err != nil, err}
}

// ----- UserStoryUsecase: Linked images -----

func (e *UserStoryUsecase) UpdateImages(id uuid.UUID, userInfo dto.BaseUserInfo,
	filesToCreate []string, filesToDelete []string) Result {

	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your story!")}
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

// ----- UserStoryUsecase: Linked events -----

func (e *UserStoryUsecase) UpdateLinkedEvents(id uuid.UUID, userInfo dto.BaseUserInfo,
	toCreate []string, toDelete []string) Result {

	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your story!")}
	}

	for _, eventId := range toCreate {
		fullEvent, getErr := e.eventUsecase.GetFullEvent(uuid.MustParse(eventId))
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

// ----- UserStoryUsecase: Linked places -----

func (e *UserStoryUsecase) UpdateLinkedPlaces(id uuid.UUID, userInfo dto.BaseUserInfo,
	toCreate []string, toDelete []string) Result {

	story, err := e.userStoryRepo.GetStoryByID(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != story.CreatedBy.ID.String() {
		return Result{false, MakeUsecaseError("Isn't your story!")}
	}

	for _, placeId := range toCreate {
		fullPlace, getErr := e.placeUsecase.GetFullPlace(uuid.MustParse(placeId))
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
