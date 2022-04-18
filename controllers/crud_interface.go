package controllers

import "github.com/gin-gonic/gin"

type CrudController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Get(ctx *gin.Context)
}
