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
	return Result{repo.EventToEvent(event), err}
}

func (e *EventEndpoint) GetFull(id uuid.UUID) (repo.Event, error) {
	return e.eventRepo.Get(id)
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
	return Result{repo.EventToEvent(event), err}
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
