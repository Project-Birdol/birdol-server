package middlewares

import (
	"net/http"
	"time"

	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

func TokenRefresh() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		if time.Since(token_info.TokenUpdated).Seconds() > 604800 - 300 {
			ctx.JSON(http.StatusContinue, gin.H {
				"result": "need_refresh",
			})
			ctx.Abort()
		}

		ctx.Next()
	}
}