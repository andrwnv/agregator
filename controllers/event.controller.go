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

func (c *EventController) Create(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if extractErr {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	var createDto dto.CreateEvent
	if err := ctx.BindJSON(&createDto); err != nil {
		misc.IncorrectRequestBodyResponse(ctx)
		return
	}

	user, err := c.userRepo.GetByEmail(payload.Email)
	if err != nil {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	region, err := c.regionRepo.GetByRegionID(createDto.RegionID)
	if err != nil {
		misc.IncorrectRequestBodyResponse(ctx)
		return
	}

	event, err := c.eventRepo.Create(createDto, user, region)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cant create event, try later.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": repo.EventToEvent(event),
	})
}

func (c *EventController) Get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := uuid.Parse(param)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Look like you attacking me",
		})
		return
	}

	event, err := c.eventRepo.Get(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": repo.EventToEvent(event),
	})
}

func (c *EventController) Update(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func (c *EventController) Delete(ctx *gin.Context) {
	param := ctx.Param("event_id")
	eventId, err := uuid.Parse(param)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Look like you attacking me",
		})
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if extractErr {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	event, err := c.eventRepo.Get(eventId)
	if err != nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	if payload.ID != event.CreatedBy.ID.String() {
		ctx.Status(http.StatusForbidden)
		return
	}

	deleteErr := c.eventRepo.Delete(eventId)
	if deleteErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cant delete event, try later.",
		})
		return
	}
	ctx.Status(http.StatusOK)
}
