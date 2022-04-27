package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controllers struct {
	userController *controllers.UserController
	authController *controllers.AuthController
	fileController *controllers.FileController
}

func NewController(
	userCtrl *controllers.UserController,
	authCtrl *controllers.AuthController,
	fileCtrl *controllers.FileController) Controllers {
	return Controllers{
		userController: userCtrl,
		authController: authCtrl,
		fileController: fileCtrl,
	}
}

func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello.")
}

// TODO: Make main router. Router != controller, controller is deps to router.

func (c *Controllers) MakeRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.CORSMiddleware())
	{
		apiV1.GET("/say_hello", SayHello)

		userGroup := apiV1.Group("/user")
		{
			userGroup.POST("/create", c.userController.Create)
			userGroup.DELETE("/delete", middleware.AuthorizeJWTMiddleware(), c.userController.Delete)
			userGroup.GET("/me", middleware.AuthorizeJWTMiddleware(), c.userController.Get)
			userGroup.PATCH("/update", middleware.AuthorizeJWTMiddleware(), c.fileController.UploadAvatarMiddleware(), c.userController.Update)
		}

		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", c.authController.Login)
		}

		fileGroup := apiV1.Group("/file")
		{
			fileGroup.GET("/receive/:filename", c.fileController.GetImage)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, "Not found")
	})

	r.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(405, "Not allowed")
	})

	return r
}
