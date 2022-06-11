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
)

type UserStoryController struct {
	usecase        *usecases.UserStoryUsecase
	fileController *FileController
}

func NewUserStoryController(storyUsecase *usecases.UserStoryUsecase, fileCtrl *FileController) *UserStoryController {
	return &UserStoryController{
		usecase:        storyUsecase,
		fileController: fileCtrl,
	}
}

func (c *UserStoryController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/user_story")
	{
		//group.GET("/") // query required
		group.GET("/:id", c.get)
		group.POST("/create", middleware.AuthorizeJWTMiddleware(), c.create)
		group.DELETE("/delete/:id", middleware.AuthorizeJWTMiddleware(), c.delete)
		group.PATCH("/update/:id", middleware.AuthorizeJWTMiddleware(), c.update)

		group.PATCH("/add_photos/:id", middleware.AuthorizeJWTMiddleware(), c.fileController.UploadImagesMiddleware(), c.createStoryImages)
		group.PATCH("/delete_photos/:id", middleware.AuthorizeJWTMiddleware(), c.deleteStoryImages)
	}
}

func (c *UserStoryController) get(ctx *gin.Context) {
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

func (c *UserStoryController) delete(ctx *gin.Context) {
	storyId, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me.") {
		return
	}

	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Delete(storyId, payload).Error, http.StatusInternalServerError) {
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UserStoryController) update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	story, err := c.usecase.GetFullStory(uuid.MustParse(ctx.Param("id")))
	if misc.HandleError(ctx, err, http.StatusNotFound) {
		return
	}

	updateDto := repo.StoryToUpdateStory(story)
	if misc.HandleError(ctx, ctx.BindJSON(&updateDto), http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Update(story.ID, updateDto, payload).Error, http.StatusForbidden) {
		return
	}

	c.usecase.UpdateLinkedPlaces(story.ID, payload, updateDto.PlaceToCreate, updateDto.PlaceToDelete)
	c.usecase.UpdateLinkedEvents(story.ID, payload, updateDto.EventToCreate, updateDto.EventToDelete)

	ctx.Status(http.StatusNoContent)
}

func (c *UserStoryController) create(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	type ExtendedCreateDto struct {
		Places  []uuid.UUID            `json:"places"`
		Events  []uuid.UUID            `json:"events"`
		Content dto.CreateUserStoryDto `json:"content"`
	}

	var createDto ExtendedCreateDto
	if misc.HandleError(ctx, ctx.BindJSON(&createDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	result := c.usecase.Create(createDto.Content, createDto.Events, createDto.Places, payload)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": result.Value,
	})
}

func (c *UserStoryController) createStoryImages(ctx *gin.Context) {
	eventId, err := uuid.Parse(ctx.Param("id"))
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

	result := c.usecase.UpdateImages(eventId, payload, loadedFiles, []string{})
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UserStoryController) deleteStoryImages(ctx *gin.Context) {
	type FilesUrl struct {
		Urls []string `json:"urls"`
	}

	eventId, err := uuid.Parse(ctx.Param("id"))
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

	result := c.usecase.UpdateImages(eventId, payload, []string{}, files.Urls)
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusNoContent)
}
