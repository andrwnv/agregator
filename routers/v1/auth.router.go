package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authController *controllers.AuthController
}

func MakeAuthRouter(authCtrl *controllers.AuthController) *AuthRouter {
	return &AuthRouter{
		authController: authCtrl,
	}
}

func (router *AuthRouter) Make(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/auth")
	{
		group.POST("/login", router.authController.Login)
	}
}
