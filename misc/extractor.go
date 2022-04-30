package misc

import (
	"encoding/json"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/gin-gonic/gin"
)

func ExtractJwtPayload(ctx *gin.Context) (user dto.BaseUserInfo, err error) {
	claims, ok := ctx.Get("token-claims")
	if !ok {
		return user, &JwtError{}
	}

	j, _ := json.Marshal(claims.(map[string]interface{}))
	user = dto.BaseUserInfo{}
	_ = json.Unmarshal(j, &user)

	return user, nil
}
