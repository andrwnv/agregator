package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/gin-gonic/gin"
)

type EventRouter struct {
	eventController *controllers.EventController
}

func MakeEventRouter(eventCtrl *controllers.EventController) *EventRouter {
	return &EventRouter{
		eventController: eventCtrl,
	}
}

func (router *EventRouter) Make(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/event")
	{
		group.GET("/:id", router.eventController.Get)
		group.POST("/create", middleware.AuthorizeJWTMiddleware(), router.eventController.Create)
		group.PATCH("/update", middleware.AuthorizeJWTMiddleware(), router.eventController.Update)
		group.DELETE("/delete", middleware.AuthorizeJWTMiddleware(), router.eventController.Delete)
	}
}
