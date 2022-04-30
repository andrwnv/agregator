package controllers

import (
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	repo *repo.UserRepo
}

func NewAuthController(r *repo.UserRepo) *AuthController {
	return &AuthController{repo: r}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var credential dto.LoginCredentials
	err := ctx.ShouldBind(&credential)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect request body!",
		})
		return
	}

	user, err := c.repo.GetByEmail(credential.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid auth info",
		})
		return
	}

	info := services.LoginInfo{
		Email:    user.Email,
		Password: user.Password,
	}

	isUserAuthenticated := services.Login(credential, info)
	if isUserAuthenticated {
		token := core.SERVER.JwtService.GenerateToken(credential.Email, repo.UserToBaseUser(user))
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Invalid auth info",
	})
}
