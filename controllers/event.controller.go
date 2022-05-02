package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/endpoints"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type EventController struct {
	endpoint *endpoints.EventEndpoint
}

func NewEventController(endpoint *endpoints.EventEndpoint) *EventController {
	return &EventController{
		endpoint: endpoint,
	}
}

func (c *EventController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/event")
	{
		group.GET("/:id", c.get)
		group.POST("/create", middleware.AuthorizeJWTMiddleware(), c.create)
		group.PATCH("/update/:event_id", middleware.AuthorizeJWTMiddleware(), c.update)
		group.DELETE("/delete/:event_id", middleware.AuthorizeJWTMiddleware(), c.delete)
	}
}

func (c *EventController) get(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	result := c.endpoint.Get(id)
	if misc.HandleError(ctx, result.Error, http.StatusNotFound) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *EventController) create(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var createDto dto.CreateEvent
	if misc.HandleError(ctx, ctx.BindJSON(&createDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	result := c.endpoint.Create(createDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": result.Value,
	})
}

func (c *EventController) update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	event, err := c.endpoint.GetFull(uuid.MustParse(ctx.Param("event_id")))
	if misc.HandleError(ctx, err, http.StatusNotFound) {
		return
	}

	updateDto := repo.EventToUpdateEvent(event)
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.endpoint.Update(event.ID, updateDto, payload).Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *EventController) delete(ctx *gin.Context) {
	eventId, err := uuid.Parse(ctx.Param("event_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.endpoint.Delete(eventId, payload).Error, http.StatusInternalServerError) {
		return
	}
	ctx.Status(http.StatusOK)
}
