package v1

import (
	"github.com/andrwnv/event-aggregator/controllers"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/gin-gonic/gin"
)

type FileRouter struct {
	fileController *controllers.FileController
}

func MakeFileRouter(fileCtrl *controllers.FileController) *FileRouter {
	return &FileRouter{
		fileController: fileCtrl,
	}
}

func (router *FileRouter) Make(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/file")
	{
		group.PATCH("/update_avatar", middleware.AuthorizeJWTMiddleware(), router.fileController.UploadAvatar)
		group.GET("/img/:filename", router.fileController.GetImage)
	}
}
