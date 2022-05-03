package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/endpoints"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type CommentController struct {
	eventEndpoint *endpoints.EventEndpoint
	placeEndpoint *endpoints.PlaceEndpoint
}

func NewCommentController(eventEndpoint *endpoints.EventEndpoint, placeEndpoint *endpoints.PlaceEndpoint) *CommentController {
	return &CommentController{
		eventEndpoint: eventEndpoint,
		placeEndpoint: placeEndpoint,
	}
}

func (c *CommentController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/comments")
	{
		eventGroup := group.Group("/event")
		{
			eventGroup.GET("/:event_id/:page/:count", c.getEventComments)
			eventGroup.POST("/create", middleware.AuthorizeJWTMiddleware(), c.createEventComment)
			eventGroup.DELETE("/delete/:id", middleware.AuthorizeJWTMiddleware(), c.deleteEventComment)
			eventGroup.PATCH("/update/:id", middleware.AuthorizeJWTMiddleware(), c.updateEventComment)
		}

		placeGroup := group.Group("/place")
		{
			placeGroup.GET("/:place_id/:page/:count", c.getPlaceComments)
			placeGroup.POST("/create", middleware.AuthorizeJWTMiddleware(), c.createPlaceComment)
			placeGroup.DELETE("/delete/:id", middleware.AuthorizeJWTMiddleware(), c.deletePlaceComment)
			placeGroup.PATCH("/update/:id", middleware.AuthorizeJWTMiddleware(), c.updatePlaceComment)
		}
	}
}

// ----- Request context processing: Event comments -----

func (c *CommentController) getEventComments(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("event_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	pageNum, _ := strconv.Atoi(ctx.Param("page"))
	count, _ := strconv.Atoi(ctx.Param("count"))

	result := c.eventEndpoint.GetComments(id, pageNum, count)
	if misc.HandleError(ctx, result.Error, http.StatusNoContent) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *CommentController) createEventComment(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var createDto dto.CreateEventCommentDto
	err := ctx.BindJSON(&createDto)
	if misc.HandleError(ctx, err, http.StatusBadRequest, "Incorrect body for comment create.") {
		return
	}

	result := c.eventEndpoint.CreateComment(createDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": result.Value,
	})
}

func (c *CommentController) deleteEventComment(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	result := c.eventEndpoint.DeleteComment(id, payload)
	if misc.HandleError(ctx, result.Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *CommentController) updateEventComment(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	var updateDto dto.UpdateEventCommentDto
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	result := c.eventEndpoint.UpdateComment(id, updateDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ----- Request context processing: Place comments -----

func (c *CommentController) getPlaceComments(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("place_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	pageNum, _ := strconv.Atoi(ctx.Param("page"))
	count, _ := strconv.Atoi(ctx.Param("count"))

	result := c.placeEndpoint.GetComments(id, pageNum, count)
	if misc.HandleError(ctx, result.Error, http.StatusNoContent) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *CommentController) createPlaceComment(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var createDto dto.CreatePlaceCommentDto
	err := ctx.BindJSON(&createDto)
	if misc.HandleError(ctx, err, http.StatusBadRequest, "Incorrect body for comment create.") {
		return
	}

	result := c.placeEndpoint.CreateComment(createDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": result.Value,
	})
}

func (c *CommentController) deletePlaceComment(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	result := c.placeEndpoint.DeleteComment(id, payload)
	if misc.HandleError(ctx, result.Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *CommentController) updatePlaceComment(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	var updateDto dto.UpdatePlaceCommentDto
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	result := c.placeEndpoint.UpdateComment(id, updateDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusNoContent)
}
