package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

// to read session_id
type ExtructSession struct {
	SessionID string `json:"session_id" binding:"required"`
}

func ReadSessionIDfromQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		access_token := token_info.Token
		session_id := ctx.Query("session_id")

		if !sessionValidityCheck(session_id, access_token) {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "error",
				"error": "login_needed",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ReadSessionIDfromBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)
		access_token := token_info.Token

		var extructor ExtructSession
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &extructor); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "error",
				"error": "something went wrong",
			})
			ctx.Abort()
			return
		}
		session_id := extructor.SessionID
		
		if !sessionValidityCheck(session_id, access_token) {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "error",
				"error": "login_needed",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func sessionValidityCheck(session_id string, access_token string) bool {
	var session model.Session
	if err := database.Sqldb.Where("session_id = ? AND access_token = ?", session_id, access_token).First(&session).Error; err != nil {
		return false
	}

	if session.Expired {
		return false
	}

	return true
}