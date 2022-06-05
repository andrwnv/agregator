package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/usecases"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type LikeController struct {
	usecase *usecases.LikeUsecase
}

func NewLikeController(usecase *usecases.LikeUsecase) *LikeController {
	return &LikeController{
		usecase: usecase,
	}
}

func (c *LikeController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/likes")
	{
		group.GET("/", middleware.AuthorizeJWTMiddleware(), c.get)
		group.GET("/is_liked/:id", middleware.AuthorizeJWTMiddleware(), c.isLiked)
		group.POST("/like", middleware.AuthorizeJWTMiddleware(), c.like)
		group.DELETE("/dislike/:id", middleware.AuthorizeJWTMiddleware(), c.dislike)
	}
}

func (c *LikeController) get(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	page, pageExtractErr := strconv.Atoi(ctx.Query("page"))
	count, countExtractErr := strconv.Atoi(ctx.Query("count"))
	if pageExtractErr != nil || countExtractErr != nil {
		page = 0
		count = 10
	}

	result := c.usecase.Get(payload, page, count)
	if misc.HandleError(ctx, result.Error, http.StatusNoContent) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *LikeController) isLiked(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	result := c.usecase.IsLiked(payload, id)
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (c *LikeController) like(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var likeDto dto.LikeDto
	if misc.HandleError(ctx, ctx.ShouldBindJSON(&likeDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	result := c.usecase.Like(likeDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": result.Value,
	})
}

func (c *LikeController) dislike(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Dislike(id, payload).Error, http.StatusInternalServerError) {
		return
	}
	ctx.Status(http.StatusOK)
}
