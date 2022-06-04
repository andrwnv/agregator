package controllers

import (
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegionController struct {
	repo *repo.RegionRepo
}

func NewRegionController(repo *repo.RegionRepo) *RegionController {
	return &RegionController{
		repo: repo,
	}
}

func (c *RegionController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/region")
	{
		group.POST("/create", c.create)
	}
}

func (c *RegionController) create(ctx *gin.Context) {
	var createDto repo.RegionDto
	if misc.HandleError(ctx, ctx.BindJSON(&createDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	c.repo.CreateRegion(createDto)
	ctx.Status(http.StatusOK)
}
