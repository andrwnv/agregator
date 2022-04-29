package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController *controllers.UserController
}

func MakeUserRouter(ctrl *controllers.UserController) *UserRouter {
	return &UserRouter{userController: ctrl}
}

func (router *UserRouter) Make(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/user")
	{
		group.POST("/create", router.userController.Create)
		group.GET("/me", middleware.AuthorizeJWTMiddleware(), router.userController.Get)
		group.DELETE("/delete", middleware.AuthorizeJWTMiddleware(), router.userController.Delete)
		group.PATCH("/update", middleware.AuthorizeJWTMiddleware(), router.userController.Update)
		group.GET("/verify/:id", router.userController.Verify)
	}
}
