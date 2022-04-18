package core

import (
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Database   *gorm.DB
	Router     *gin.Engine
	JwtService services.JWTService
}

var ServerInst *Server
