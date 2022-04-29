package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	repo *repo.UserRepo
}

func NewUserController(r *repo.UserRepo) *UserController {
	return &UserController{repo: r}
}

func (c *UserController) Create(ctx *gin.Context) {
	var _dto dto.CreateUser
	if err := ctx.BindJSON(&_dto); err != nil {
		misc.IncorrectRequestBodyResponse(ctx)
		return
	}

	if err := c.repo.Create(_dto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "User already exists!",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": "Successful create, Welcome!",
	})
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
