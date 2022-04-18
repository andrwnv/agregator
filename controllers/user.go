package controllers

import (
	"encoding/json"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
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
	claims, ok := ctx.Get("token-claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cant extract info from claims",
		})
		return
	}

	j, _ := json.Marshal(claims.(map[string]interface{}))
	user := dto.BaseUserInfo{}
	_ = json.Unmarshal(j, &user)

	if c.repo.Delete(user) != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) Get(ctx *gin.Context) {
	claims, ok := ctx.Get("token-claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cant extract info from claims",
		})
		return
	}

	j, _ := json.Marshal(claims.(map[string]interface{}))
	user := dto.BaseUserInfo{}
	_ = json.Unmarshal(j, &user)

	ctx.JSON(http.StatusCreated, gin.H{
		"result": user,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	// TODO: impl

	ctx.JSON(http.StatusOK, gin.H{
		"result": "TODO",
	})
}
