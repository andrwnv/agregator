package controllers

import (
	"github.com/andrwnv/event-aggregator/core"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/models"
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(ctx *gin.Context) {
	var credential dto.LoginCredentials
	err := ctx.ShouldBind(&credential)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect request body!",
		})
		return
	}

	user := models.User{}
	models.GetByEmail(&user, credential.Email)

	info := services.LoginInfo{
		Email:    user.Email,
		Password: user.Password,
	}

	isUserAuthenticated := services.Login(credential, info)
	if isUserAuthenticated {
		token := core.ServerInst.JwtService.GenerateToken(credential.Email, models.To(user))
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid auth info",
			})
		}
	}
}
