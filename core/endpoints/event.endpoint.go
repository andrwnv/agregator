package endpoints

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/google/uuid"
)

type EventEndpoint struct {
	eventRepo    *repo.EventRepo
	userEndpoint *UserEndpoint
	regionRepo   *repo.RegionRepo
}

func NewEventEndpoint(
	eventRepo *repo.EventRepo,
	userEndpoint *UserEndpoint,
	regionRepo *repo.RegionRepo) *EventEndpoint {
	return &EventEndpoint{
		eventRepo:    eventRepo,
		userEndpoint: userEndpoint,
		regionRepo:   regionRepo,
	}
}

func (e *EventEndpoint) Get(id uuid.UUID) Result {
	event, err := e.eventRepo.Get(id)
	eventPhotos, _ := e.eventRepo.GetImages(id)
	return Result{repo.EventToEvent(event, eventPhotos), err}
}

func (e *EventEndpoint) GetFullEvent(id uuid.UUID) (repo.Event, error) {
	return e.eventRepo.Get(id)
}

func (e *EventEndpoint) GetFullEventComment(id uuid.UUID) (repo.EventComment, error) {
	return e.eventRepo.GetCommentByID(id)
}

func (e *EventEndpoint) Create(createDto dto.CreateEvent, userInfo dto.BaseUserInfo) Result {
	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find user for create event.")}
	}

	region, err := e.regionRepo.GetByRegionID(createDto.RegionID)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find selected country.")}
	}

	// TODO: check begin, end datetime correctness for upd & create

	event, err := e.eventRepo.Create(createDto, user, region)
	return Result{repo.EventToEvent(event, []string{}), err}
}

func (e *EventEndpoint) Update(id uuid.UUID, updateDto dto.UpdateEvent, userInfo dto.BaseUserInfo) Result {
	event, err := e.eventRepo.Get(id)
	if err != nil {
		return Result{nil, err}
	}

	if userInfo.ID != event.CreatedBy.ID.String() {
		return Result{nil, MakeEndpointError("Isn't your event!")}
	}

	event.Region, err = e.regionRepo.GetByRegionID(updateDto.RegionID)
	if err != nil {
		return Result{nil, MakeEndpointError("Cant find selected country.")}
	}

	// TODO: check begin, end datetime correctness for upd & create

	err = e.eventRepo.Update(event.ID, updateDto, event.Region)
	return Result{err != nil, err}
}

func (e *EventEndpoint) Delete(id uuid.UUID, userInfo dto.BaseUserInfo) Result {
	event, err := e.eventRepo.Get(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != event.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your event!")}
	}

	err = e.eventRepo.Delete(id)
	return Result{err != nil, err}
}

// ----- EventEndpoint: Images -----

func (e *EventEndpoint) UpdateEventImages(id uuid.UUID, userInfo dto.BaseUserInfo,
	filesToCreate []string, filesToDelete []string) Result {

	event, err := e.eventRepo.Get(id)
	if err != nil {
		return Result{false, err}
	}

	if userInfo.ID != event.CreatedBy.ID.String() {
		return Result{false, MakeEndpointError("Isn't your event!")}
	}

	for _, url := range filesToCreate {
		err := e.eventRepo.CreateImages(event.ID, url)
		if err != nil {
			return Result{false, err}
		}
	}

	for _, url := range filesToDelete {
		err := e.eventRepo.DeleteImages(url)
		// TODO: delete photos from dir.
		if err != nil {
			return Result{false, err}
		}
	}

	return Result{true, nil}
}

// ----- EventEndpoint: Comments -----

func (e *EventEndpoint) CreateComment(createDto dto.CreateEventCommentDto, userInfo dto.BaseUserInfo) Result {
	user, err := e.userEndpoint.GetFull(userInfo)
	if err != nil {
		return Result{nil, err}
	}
	event, err := e.eventRepo.Get(uuid.MustParse(createDto.LinkedEventID))
	if err != nil {
		return Result{false, err}
	}

	comment, err := e.eventRepo.CreateComment(createDto, user, event)
	if err != nil {
		return Result{nil, MakeEndpointError("Failed to create comment.")}
	}

	return Result{repo.CommentToComment(comment), nil}
}

func (e *EventEndpoint) GetComments(eventId uuid.UUID, page int, count int) Result {
	comments, err := e.eventRepo.GetComments(eventId, page, count)
	if err != nil {
		return Result{nil, MakeEndpointError("Failed to create comment.")}
	}

	var result []dto.EventCommentDto
	for _, value := range comments {
		result = append(result, repo.CommentToComment(value))
	}

	return Result{result, nil}
}

func (e *EventEndpoint) DeleteComment(commentId uuid.UUID, userInfo dto.BaseUserInfo) Result {
	comment, err := e.eventRepo.GetCommentByID(commentId)
	if err != nil {
		return Result{nil, err}
	}
	if userInfo.ID != comment.CreatedBy.ID.String() {
		return Result{nil, MakeEndpointError("Isn't your comment!")}
	}

	err = e.eventRepo.DeleteComments(commentId)
	if err != nil {
		return Result{false, MakeEndpointError("Cant delete comment(s).")}
	}
	return Result{true, nil}
}

func (e *EventEndpoint) UpdateComment(id uuid.UUID, updateDto dto.UpdateEventCommentDto, userInfo dto.BaseUserInfo) Result {
	comment, err := e.eventRepo.GetCommentByID(id)
	if err != nil {
		return Result{nil, err}
	}
	if userInfo.ID != comment.CreatedBy.ID.String() {
		return Result{nil, MakeEndpointError("Isn't your comment!")}
	}

	err = e.eventRepo.UpdateComment(id, updateDto)
	if err != nil {
		return Result{false, MakeEndpointError("Cant update comment(s).")}
	}
	return Result{true, nil}
}
