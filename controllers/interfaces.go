package controllers

import "github.com/gin-gonic/gin"

type CrudController interface {
	create(ctx *gin.Context)
	update(ctx *gin.Context)
	delete(ctx *gin.Context)
	get(ctx *gin.Context)
}

type IRouter interface {
	MakeRoutesV1(rootGroup *gin.RouterGroup)
}
