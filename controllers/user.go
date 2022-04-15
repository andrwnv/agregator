package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func RegisterUser(ctx *gin.Context) {
	var _dto dto.CreateUser

	if err := ctx.BindJSON(&_dto); err != nil {
		ctx.JSON(http.StatusBadRequest, "Incorrect req body")
	}

	user := models.User{
		ID:        uuid.New(),
		FirstName: _dto.FirstName,
		LastName:  _dto.SecondName,
		Email:     _dto.Email,
		Password:  _dto.Password,
	}
	err := models.CreateUser(&user)

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "User already exists!",
		})
		return
	}

	ctx.JSON(http.StatusCreated, &dto.BaseUserInfo{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		SecondName: user.LastName,
		Email:      user.Email,
	})
}
