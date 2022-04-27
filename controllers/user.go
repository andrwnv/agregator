package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/utils"
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect request body!",
		})
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
	user, err := utils.ExtractJwtPayload(ctx)
	if err {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cant extract info from claims",
		})
		return
	}

	if c.repo.Delete(user) != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) Get(ctx *gin.Context) {
	user, err := utils.ExtractJwtPayload(ctx)
	if err {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cant extract info from claims",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": user,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	val := ctx.GetString("file-name")
	utils.ReportInfo(val)

	ctx.JSON(http.StatusOK, gin.H{
		"result": "test",
	})
}
