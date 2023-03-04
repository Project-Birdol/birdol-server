/*
Authrization middleware for gin
*/
package middleware

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/Project-Birdol/birdol-server/database/model"
	"github.com/Project-Birdol/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
)

// to read session_id
type ExtractSession struct {
	SessionID string `json:"session_id"`
}

type SessionMiddleware struct {
	DB *gorm.DB
}

func (sm *SessionMiddleware) ReadSessionIDfromQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		access_token := token_info.Token
		session_id := ctx.Query("session_id")

		if !sm.sessionValidityCheck(session_id, access_token) {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrNotLoggedIn)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (sm *SessionMiddleware) ReadSessionIDfromBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)
		access_token := token_info.Token

		content_type := ctx.GetHeader("Content-Type")
		if content_type != gin.MIMEJSON {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidType)
			ctx.Abort()
			return
		}

		var extractor ExtractSession
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &extractor); err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailParseJSON)
			ctx.Abort()
			return
		}
		session_id := extractor.SessionID

		if !sm.sessionValidityCheck(session_id, access_token) {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrNotLoggedIn)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (sm *SessionMiddleware) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		if time.Since(token_info.TokenUpdated).Seconds() > 604800-300 {
			response.SetNormalResponse(ctx, http.StatusAccepted, response.ResultNeedTokenRefresh)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (sm *SessionMiddleware) sessionValidityCheck(session_id string, access_token string) bool {
	var session model.Session
	if err := sm.DB.Where("session_id = ? AND access_token = ?", session_id, access_token).First(&session).Error; err != nil {
		return false
	}

	return !session.Expired
}
