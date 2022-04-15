package controllers

import (
	"crypto/sha1"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUser(ctx *gin.Context) {
	var _dto dto.CreateUser

	if err := ctx.BindJSON(&_dto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect request body!",
		})
		return
	}

	passHash := sha1.New()
	passHash.Write([]byte(_dto.Password))

	user := models.From(_dto)
	err := models.CreateUser(&user)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "User already exists!",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.To(user))
}
