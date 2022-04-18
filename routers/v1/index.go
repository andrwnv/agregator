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
}

func NewController(
	u *controllers.UserController,
	a *controllers.AuthController) Controllers {
	return Controllers{
		userController: u,
		authController: a,
	}
}

func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello.")
}

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
			userGroup.GET("/me", c.userController.Get)
		}

		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", c.authController.Login)
		}
	}

	return r
}
