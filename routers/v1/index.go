package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello.")
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/say_hello", SayHello)
	}

	return r
}
