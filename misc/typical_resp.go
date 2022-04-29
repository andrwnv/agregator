package misc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FailedClaimsExtractResponse(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": "Cant extract info from claims",
	})
}

func IncorrectRequestBodyResponse(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": "Incorrect request body!",
	})
}
