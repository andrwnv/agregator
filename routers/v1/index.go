package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/gin-gonic/gin"
)

func MakeRouter(
	userCtrl *controllers.UserController,
	eventCtrl *controllers.EventController,
	placeCtrl *controllers.PlaceController,
	authCtrl *controllers.AuthController,
	fileCtrl *controllers.FileController,
	commentsCtrl *controllers.CommentController,
	userStoryCtrl *controllers.UserStoryController) *gin.Engine {

	engine := gin.New()
	engine.Use(gin.Logger())

	v1Group := engine.Group("/api/v1")
	v1Group.Use(middleware.CORSMiddleware())
	{
		userCtrl.MakeRoutesV1(v1Group)
		eventCtrl.MakeRoutesV1(v1Group)
		placeCtrl.MakeRoutesV1(v1Group)
		authCtrl.MakeRoutesV1(v1Group)
		fileCtrl.MakeRoutesV1(v1Group)
		commentsCtrl.MakeRoutesV1(v1Group)
		userStoryCtrl.MakeRoutesV1(v1Group)
	}

	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, "Not found")
	})

	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(405, "Not allowed")
	})

	return engine
}
