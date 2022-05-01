package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type EventController struct {
	eventRepo  *repo.EventRepo
	userRepo   *repo.UserRepo
	regionRepo *repo.RegionRepo
}

func NewEventController(eventRepo *repo.EventRepo,
	userRepo *repo.UserRepo,
	regionRepo *repo.RegionRepo) *EventController {

	return &EventController{
		eventRepo:  eventRepo,
		userRepo:   userRepo,
		regionRepo: regionRepo,
	}
}

// TODO: check begin, end datetime correctness for upd & create

func (c *EventController) Create(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var createDto dto.CreateEvent
	if misc.HandleError(ctx, ctx.BindJSON(&createDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	user, err := c.userRepo.GetByEmail(payload.Email)
	if misc.HandleError(ctx, err, http.StatusInternalServerError) {
		return
	}

	region, err := c.regionRepo.GetByRegionID(createDto.RegionID)
	if misc.HandleError(ctx, err, http.StatusBadRequest, "Cant find selected country.") {
		return
	}

	event, err := c.eventRepo.Create(createDto, user, region)
	if misc.HandleError(ctx, err, http.StatusInternalServerError, "Cant create event, try later.") {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": repo.EventToEvent(event),
	})
}

func (c *EventController) Get(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	event, err := c.eventRepo.Get(id)
	if misc.HandleError(ctx, err, http.StatusNotFound, "Event not found") {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": repo.EventToEvent(event),
	})
}

func (c *EventController) Update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	event, err := c.eventRepo.Get(uuid.MustParse(ctx.Param("event_id")))
	if misc.HandleError(ctx, err, http.StatusNotFound) {
		return
	}

	updateDto := repo.EventToUpdateEvent(event)
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	if payload.ID != event.CreatedBy.ID.String() {
		ctx.Status(http.StatusForbidden)
		return
	}

	event.Region, err = c.regionRepo.GetByRegionID(updateDto.RegionID)
	if misc.HandleError(ctx, err, http.StatusBadRequest, "Cant find selected country.") {
		return
	}

	if misc.HandleError(ctx, c.eventRepo.Update(event.ID, updateDto, event.Region),
		http.StatusInternalServerError, "Cant update event, try later.") {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *EventController) Delete(ctx *gin.Context) {
	eventId, err := uuid.Parse(ctx.Param("event_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	event, err := c.eventRepo.Get(eventId)
	if misc.HandleError(ctx, err, http.StatusNotFound) {
		return
	}

	if payload.ID != event.CreatedBy.ID.String() {
		ctx.Status(http.StatusForbidden)
		return
	}

	if misc.HandleError(ctx, c.eventRepo.Delete(eventId), http.StatusInternalServerError, "Cant delete event, try later.") {
		return
	}

	ctx.Status(http.StatusOK)
}
