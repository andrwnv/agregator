package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/usecases"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserController struct {
	usecase *usecases.UserUsecase
}

func NewUserController(usecase *usecases.UserUsecase) *UserController {
	return &UserController{
		usecase: usecase,
	}
}

func (c *UserController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/user")
	{
		group.POST("/create", c.create)
		group.GET("/me", middleware.AuthorizeJWTMiddleware(), c.get)
		group.DELETE("/delete", middleware.AuthorizeJWTMiddleware(), c.delete)
		group.GET("/:id", c.getByID)
		group.PATCH("/update", middleware.AuthorizeJWTMiddleware(), c.update)
		group.GET("/verify/:id", c.verify)
	}
}

// ----- Request context processing -----

func (c *UserController) get(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": payload,
	})
}

func (c *UserController) create(ctx *gin.Context) {
	var createDto dto.CreateUser
	if misc.HandleError(ctx, ctx.BindJSON(&createDto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	res := c.usecase.Create(createDto)
	if misc.HandleError(ctx, res.Error, http.StatusConflict) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": res.Value,
	})
}

func (c *UserController) delete(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.usecase.Delete(payload).Error, http.StatusInternalServerError) {
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) getByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	result := c.usecase.GetByID(id)
	if misc.HandleError(ctx, result.Error, http.StatusBadRequest) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}

func (c *UserController) update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	user, err := c.usecase.GetFull(payload)
	if misc.HandleError(ctx, err, http.StatusInternalServerError, "Cant extract user from database") {
		return
	}

	updateDto := repo.UserToUpdateUserDto(user)
	_ = ctx.BindJSON(&updateDto)

	result := c.usecase.Update(user.ID, updateDto)
	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError) {
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"result": result.Value,
	})
}

func (c *UserController) verify(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me") {
		return
	}

	result := c.usecase.Verify(id)
	if misc.HandleError(ctx, result.Error, http.StatusForbidden) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}
