package core

import (
	"github.com/andrwnv/event-aggregator/core/services"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	JwtService services.JWTService
}

var SERVER *Server
