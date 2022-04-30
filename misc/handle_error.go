package misc

import (
	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, err error, status int, text ...string) bool {
	haveError := err != nil
	if haveError {
		if len(text) == 0 {
			ctx.JSON(status, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(status, gin.H{
				"error": text[0],
			})
		}
	}
	return haveError
}
