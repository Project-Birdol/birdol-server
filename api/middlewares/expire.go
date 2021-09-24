package middlewares

import (
	"net/http"
	"time"

	"github.com/MISW/birdol-server/database/model"
	"github.com/MISW/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
)

func CheckToken() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		if time.Since(token_info.TokenUpdated).Seconds() > 604800 - 300 {
			response.SetNormalResponse(ctx, http.StatusAccepted, response.ResultNeedTokenRefresh)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}