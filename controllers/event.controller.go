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

	event, err := c.eventRepo.Create(createDto, user)
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
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": repo.EventToEvent(event),
	})
}

func (c *EventController) Update(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func (c *EventController) Delete(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
