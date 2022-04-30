package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserController struct {
	repo   *repo.UserRepo
	mailer *services.Mailer
}

func NewUserController(r *repo.UserRepo, mailer *services.Mailer) *UserController {
	return &UserController{
		repo:   r,
		mailer: mailer,
	}
}

func (c *UserController) Create(ctx *gin.Context) {
	var _dto dto.CreateUser
	if misc.HandleError(ctx, ctx.BindJSON(&_dto), http.StatusBadRequest, "Incorrect request body.") {
		return
	}

	user, err := c.repo.Create(_dto)
	if misc.HandleError(ctx, err, http.StatusConflict, "User already exists!") {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": "Successful create, Welcome!",
	})

	go func() {
		to := []string{user.Email}
		err := c.mailer.SendVerifyEmail(to, user.ID.String())
		if err != nil {
			misc.ReportError("Cant sent verify email!")
		}
	}()
}

func (c *UserController) Delete(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	if misc.HandleError(ctx, c.repo.Delete(payload), http.StatusInternalServerError, "Cant delete user, try later.") {
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) Get(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": payload,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, extractErr, http.StatusBadRequest) {
		return
	}

	user, err := c.repo.GetByEmail(payload.Email)
	if misc.HandleError(ctx, err, http.StatusInternalServerError, "Cant extract user from database") {
		return
	}

	updateDto := repo.UserToUpdateUserDto(user)
	_ = ctx.BindJSON(&updateDto)

	result, err := c.repo.Update(user.ID, updateDto)
	if misc.HandleError(ctx, err, http.StatusInternalServerError, "Cant update user info, try later.") {
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"result": result,
	})
}

func (c *UserController) Verify(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if misc.HandleError(ctx, err, http.StatusForbidden, "Look like you attacking me") {
		return
	} else if misc.HandleError(ctx, c.repo.Verify(id), http.StatusForbidden, "Something went wrong. Try later.") {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "User verified",
	})
}
