package controllers

import (
	"errors"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/usecases"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type PlaceController struct {
	usecase        *usecases.PlaceUsecase
	fileController *FileController
}

func NewPlaceController(usecase *usecases.PlaceUsecase, fileCtrl *FileController) *PlaceController {
	return &PlaceController{
		usecase:        usecase,
		fileController: fileCtrl,
	}
}

func (c *PlaceController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/place")
	{
		group.GET("/:id", c.get)
		group.GET("/consume/:page/:count", c.getPlaces)
		group.POST("/create", middleware.AuthorizeJWTMiddleware(), c.create)
		group.PATCH("/update/:place_id", middleware.AuthorizeJWTMiddleware(), c.update)
		group.DELETE("/delete/:place_id", middleware.AuthorizeJWTMiddleware(), c.delete)

		group.PATCH("/add_photos/:place_id", middleware.AuthorizeJWTMiddleware(), c.fileController.UploadImagesMiddleware(), c.createPlaceImages)
		group.PATCH("/delete_photos/:place_id", middleware.AuthorizeJWTMiddleware(), c.deletePlaceImages)
	}
}

func (c *PlaceController) get(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	result := c.usecase.Get(id)
	if misc.HandleError(ctx, result.Error, http.StatusNotFound) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *PlaceController) getPlaces(ctx *gin.Context) {
	pageNum, _ := strconv.Atoi(ctx.Param("page"))
	count, _ := strconv.Atoi(ctx.Param("count"))

	result := c.usecase.GetPlaces(pageNum, count)
	if misc.HandleError(ctx, result.Error, http.StatusNoContent) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *PlaceController) create(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var createDto dto.CreatePlace
	if misc.HandleError(ctx, ctx.BindJSON(&createDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	result := c.usecase.Create(createDto, payload)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": result.Value,
	})
}

func (c *PlaceController) update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	place, err := c.usecase.GetFullPlace(uuid.MustParse(ctx.Param("place_id")))
	if misc.HandleError(ctx, err, http.StatusNotFound) {
		return
	}

	updateDto := repo.PlaceToUpdatePlace(place)
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Update(place.ID, updateDto, payload).Error, http.StatusForbidden) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PlaceController) delete(ctx *gin.Context) {
	placeId, err := uuid.Parse(ctx.Param("place_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Delete(placeId, payload).Error, http.StatusInternalServerError) {
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *PlaceController) createPlaceImages(ctx *gin.Context) {
	placeId, err := uuid.Parse(ctx.Param("place_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	loadedFiles := ctx.GetStringSlice("file-names")
	if len(loadedFiles) == 0 {
		if misc.HandleError(ctx, errors.New("no images loaded"), http.StatusBadRequest) {
			return
		}
	}

	result := c.usecase.UpdatePlaceImages(placeId, payload, loadedFiles, []string{})
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PlaceController) deletePlaceImages(ctx *gin.Context) {
	type FilesUrl struct {
		Urls []string `json:"urls"`
	}

	placeId, err := uuid.Parse(ctx.Param("place_id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	var files FilesUrl
	if misc.HandleError(ctx, ctx.BindJSON(&files), http.StatusBadRequest, "No file to delete.") {
		return
	}

	result := c.usecase.UpdatePlaceImages(placeId, payload, []string{}, files.Urls)
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusNoContent)
}
