package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/endpoints"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"net/http"
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
		group.GET("/event/:event_id/:page/:count", c.get)
		group.POST("/event/create", middleware.AuthorizeJWTMiddleware(), c.create)
		group.DELETE("/event/delete/:id", c.delete)
		group.PATCH("/event/update/:id", c.update)
	}
}

// ----- Request context processing -----

func (c *CommentController) get(ctx *gin.Context) {

}

func (c *CommentController) create(ctx *gin.Context) {
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

func (c *CommentController) delete(ctx *gin.Context) {

}

func (c *CommentController) update(ctx *gin.Context) {

}
