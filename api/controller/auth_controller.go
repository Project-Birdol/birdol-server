package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Project-Birdol/birdol-server/auth"
	"github.com/Project-Birdol/birdol-server/controller/jsonmodel"
	"github.com/Project-Birdol/birdol-server/database"
	"github.com/Project-Birdol/birdol-server/database/model"
	"github.com/Project-Birdol/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
)

/*
	SetDataLink : Loginするためのpasswordを設定, 更新する
*/
func SetDataLink() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[SetDataLink] ")
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		// request_body data の jsonを変換
		var request_body jsonmodel.EnableLinkRequest
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &request_body); err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailParseJSON)
			return
		}

		// request data に含まれるパスワードをハッシュ化する
		if err := auth.HashString(&request_body.Password); err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailSetPassword)
			return
		}

		// Set expire date
		expire_day := time.Now().Add(time.Hour * 24 * 7)

		// update database
		if err := database.Sqldb.Model(&model.User{}).Where("id = ?", token_info.UserID).Updates(map[string]interface{}{"password": request_body.Password, "expire_date": expire_day}).Error; err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailSetPassword)
			return
		}

		// response
		property := gin.H { "expire_date": expire_day.Format("2006-01-02 15:04:05") }
		response.SetNormalResponse(ctx, http.StatusOK, response.ResultOK, property)
	}
}

/*
	UnlinkAccount : Logoutしてaccess_token，関連sessionを削除する
*/
func UnlinkAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[UnlinkAccount] ")
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		//request data のrequest_bodyを変換
		var request_body jsonmodel.AuthLogoutRequest
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &request_body); err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailParseJSON)
			return
		}

		access_token := token_info.Token
		device_id := token_info.DeviceID

		// logoutリクエストのため、access tokenを削除する
		if err := auth.DeleteToken(access_token, device_id); err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailUnlink)
			return
		}

		if err := database.Sqldb.Where("access_token = ?", access_token).Delete(&model.Session{}).Error; err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailUnlink)
			return
		}

		//レスポンス
		response.SetNormalResponse(ctx, http.StatusOK, response.ResultOK)
	}
}

/*
	Token Authorization Handler
*/
func TokenAuthorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[TokenAuthorize] ")
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)

		user_id := token_info.UserID
		access_token := token_info.Token
		device_id := token_info.DeviceID

		session_id, err := auth.CreateSession(device_id, access_token, user_id)
		if err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailCreateSession)
			return
		}

		property := gin.H { "session_id": session_id }
		response.SetNormalResponse(ctx, http.StatusOK, response.ResultOK, property)
	}
}

/*
  Regenerate token using refresh_token
*/
func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[RefreshToken] ")
		token_interface, _ := ctx.Get("access_token")
		token_info := token_interface.(model.AccessToken)
		refresh_token := ctx.Query("refresh_token")

		user_id := token_info.UserID
		device_id := token_info.DeviceID

		if refresh_token != token_info.RefreshToken {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidRefreshToken)
			return
		}
				
		new_token, new_refresh, err := auth.SetToken(token_info.UserID, token_info.DeviceID, token_info.PublicKey, token_info.KeyType)
		if err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailRefresh)
			return
		}

		session_id, err := auth.CreateSession(device_id, new_token, user_id)
		if err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailCreateSession)
			return
		}

		property := gin.H {
			"token": new_token,
			"refresh_token": new_refresh,
			"session_id": session_id,
		}
		response.SetNormalResponse(ctx, http.StatusOK, response.ResultRefreshSuccess, property)
	}
}