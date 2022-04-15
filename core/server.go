package core

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Database *gorm.DB
	Router   *gin.Engine
}

var ServerInst *Server
