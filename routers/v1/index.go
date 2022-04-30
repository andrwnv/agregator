package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userRouter  *UserRouter
	fileRouter  *FileRouter
	authRouter  *AuthRouter
	eventRouter *EventRouter
}

func MakeRouter(
	userCtrl *controllers.UserController,
	authCtrl *controllers.AuthController,
	fileCtrl *controllers.FileController,
	eventCtrl *controllers.EventController) *Router {

	return &Router{
		userRouter:  MakeUserRouter(userCtrl),
		fileRouter:  MakeFileRouter(fileCtrl),
		authRouter:  MakeAuthRouter(authCtrl),
		eventRouter: MakeEventRouter(eventCtrl),
	}
}

func (router *Router) InitRouter() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())

	v1Group := engine.Group("/api/v1")
	v1Group.Use(middleware.CORSMiddleware())
	{
		router.userRouter.Make(v1Group)
		router.authRouter.Make(v1Group)
		router.fileRouter.Make(v1Group)
		router.eventRouter.Make(v1Group)
	}

	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, "Not found")
	})

	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(405, "Not allowed")
	})

	return engine
}
