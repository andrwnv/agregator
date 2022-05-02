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
	endpoint *endpoints.EventEndpoint
}

func NewCommentController(endpoint *endpoints.EventEndpoint) *CommentController {
	return &CommentController{
		endpoint: endpoint,
	}
}

func (c *CommentController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/comments")
	{
		group.GET("/event/:event_id/:page/:count", c.getEventComments)
		group.POST("/event/create", middleware.AuthorizeJWTMiddleware(), c.createEventComment)
		group.DELETE("/event/delete/:id", middleware.AuthorizeJWTMiddleware(), c.deleteEventComment)
		group.PATCH("/event/update/:id", middleware.AuthorizeJWTMiddleware(), c.updateEventComment)
	}
}

// ----- Request context processing -----

func (c *CommentController) getEventComments(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("event_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	pageNum, _ := strconv.Atoi(ctx.Param("page"))
	count, _ := strconv.Atoi(ctx.Param("count"))

	result := c.endpoint.GetComments(id, pageNum, count)
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

	result := c.endpoint.CreateComment(createDto, payload)
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

	result := c.endpoint.DeleteComment(id, payload)
	if misc.HandleError(ctx, result.Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *CommentController) updateEventComment(ctx *gin.Context) {

}
