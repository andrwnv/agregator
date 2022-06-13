package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/usecases"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type SearchController struct {
	usecase *usecases.SearchUsecase
}

func NewSearchController(searchUsecase *usecases.SearchUsecase) *SearchController {
	return &SearchController{
		usecase: searchUsecase,
	}
}

func (c *SearchController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/search")
	{
		group.GET("/nearby", c.searchNearby)
		group.GET("/value", c.searchByValue)
	}
}

func (c *SearchController) searchNearby(ctx *gin.Context) {
	latStr, latOk := ctx.GetQuery("lat")
	lonStr, lonOk := ctx.GetQuery("lon")
	fromStr, fromOk := ctx.GetQuery("f")
	sizeStr, sizeOk := ctx.GetQuery("s")
	rawObjTypes, _ := ctx.GetQuery("type")

	if !latOk || !lonOk || !fromOk || !sizeOk {
		ctx.Status(http.StatusBadRequest)
		return
	}

	objTypes := strings.Split(rawObjTypes, ",")
	if objTypes[0] == "" {
		objTypes = []string{"event", "place"}
	}

	lat, err := strconv.ParseFloat(latStr, 32)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 32)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	size, err := strconv.ParseUint(sizeStr, 10, 16)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	from, err := strconv.ParseUint(fromStr, 10, 16)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	result := c.usecase.SearchNearby(dto.SearchNearbyDto{
		Coords: dto.LocationDto{
			Lat: float32(lat),
			Lon: float32(lon),
		},
		SearchType: objTypes,
		From:       uint(from),
		Limit:      uint(size),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *SearchController) searchByValue(ctx *gin.Context) {
	value, valueOk := ctx.GetQuery("val")
	fromStr, fromOk := ctx.GetQuery("f")
	sizeStr, sizeOk := ctx.GetQuery("s")
	rawObjTypes, _ := ctx.GetQuery("type")

	if !valueOk || !fromOk || !sizeOk {
		ctx.Status(http.StatusBadRequest)
		return
	}

	objTypes := strings.Split(rawObjTypes, ",")
	if objTypes[0] == "" {
		objTypes = []string{"event", "place"}
	}

	size, err := strconv.ParseUint(sizeStr, 10, 16)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	from, err := strconv.ParseUint(fromStr, 10, 16)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	result := c.usecase.Search(dto.SearchDto{
		ValueToSearch: value,
		SearchType:    objTypes,
		From:          uint(from),
		Limit:         uint(size),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}
