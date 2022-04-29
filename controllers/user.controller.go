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
	if err := ctx.BindJSON(&_dto); err != nil {
		misc.IncorrectRequestBodyResponse(ctx)
		return
	}

	user, err := c.repo.Create(_dto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "User already exists!",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": "Successful create, Welcome!",
	})

	go func() {
		to := []string{user.Email}
		err := c.mailer.SendVerifyEmail(to, user.ID.String())
		if err != nil {
			misc.ReportError("Cant sent verify email")
		}
	}()
}

func (c *UserController) Delete(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if err {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	if c.repo.Delete(payload) != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) Get(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if err {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": payload,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	payload, extractErr := misc.ExtractJwtPayload(ctx)
	if extractErr {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	user, err := c.repo.GetByEmail(payload.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Cant extract user from database",
		})
		return
	}

	updateDto := repo.ToUpdateDto(user)
	_ = ctx.BindJSON(&updateDto)

	result, err := c.repo.Update(user.ID, updateDto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Cant update user info!",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"result": result,
	})
}

func (c *UserController) Verify(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Look like you attacking me",
		})
		return
	}

	if err := c.repo.Verify(uuid); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong. Try later",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "User verified",
	})
}
