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

type EventController struct {
	usecase        *usecases.EventUsecase
	fileController *FileController
}

func NewEventController(usecase *usecases.EventUsecase, fileCtrl *FileController) *EventController {
	return &EventController{
		usecase:        usecase,
		fileController: fileCtrl,
	}
}

func (c *EventController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/event")
	{
		group.GET("/:id", c.get)
		group.GET("/consume/:page/:count", c.getEvents)
		group.POST("/create", middleware.AuthorizeJWTMiddleware(), c.create)
		group.PATCH("/update/:event_id", middleware.AuthorizeJWTMiddleware(), c.update)
		group.DELETE("/delete/:event_id", middleware.AuthorizeJWTMiddleware(), c.delete)

		group.PATCH("/add_photos/:event_id", middleware.AuthorizeJWTMiddleware(), c.fileController.UploadImagesMiddleware(), c.createEventImages)
		group.PATCH("/delete_photos/:event_id", middleware.AuthorizeJWTMiddleware(), c.deleteEventImages)
	}
}

func (c *EventController) get(ctx *gin.Context) {
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

func (c *EventController) getEvents(ctx *gin.Context) {
	pageNum, _ := strconv.Atoi(ctx.Param("page"))
	count, _ := strconv.Atoi(ctx.Param("count"))

	result := c.usecase.GetEvents(pageNum, count)
	if misc.HandleError(ctx, result.Error, http.StatusNoContent) {
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

	result := c.usecase.Create(createDto, payload)
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

	event, err := c.usecase.GetFullEvent(uuid.MustParse(ctx.Param("event_id")))
	if misc.HandleError(ctx, err, http.StatusNotFound) {
		return
	}

	updateDto := repo.EventToUpdateEvent(event)
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Update(event.ID, updateDto, payload).Error, http.StatusForbidden) {
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

	if misc.HandleError(ctx, c.usecase.Delete(eventId, payload).Error, http.StatusInternalServerError) {
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *EventController) createEventImages(ctx *gin.Context) {
	eventId, err := uuid.Parse(ctx.Param("event_id"))
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

	result := c.usecase.UpdateEventImages(eventId, payload, loadedFiles, []string{})
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *EventController) deleteEventImages(ctx *gin.Context) {
	type FilesUrl struct {
		Urls []string `json:"urls"`
	}

	eventId, err := uuid.Parse(ctx.Param("event_id"))
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

	result := c.usecase.UpdateEventImages(eventId, payload, []string{}, files.Urls)
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusNoContent)
}
