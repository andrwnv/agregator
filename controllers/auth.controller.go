package controllers

import (
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/core/endpoints"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	endpoint *endpoints.AuthEndpoint
}

func NewAuthController(endpoint *endpoints.AuthEndpoint) *AuthController {
	return &AuthController{endpoint: endpoint}
}

func (c *AuthController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/auth")
	{
		group.POST("/login", c.Login)
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var credential dto.LoginCredentials
	err := ctx.ShouldBind(&credential)
	if misc.HandleError(ctx, err, http.StatusUnauthorized, "Incorrect request body") {
		return
	}

	result := c.endpoint.Login(credential)
	if misc.HandleError(ctx, result.Error, http.StatusUnauthorized) {
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"result": result.Value,
	})
}
