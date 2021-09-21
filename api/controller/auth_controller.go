package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/MISW/birdol-server/auth"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
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

		// request data の jsonを変換
		var json jsonmodel.EnableLinkRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}

		// request data に含まれるパスワードをハッシュ化する
		if err := auth.HashString(&json.Password); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データ連携の設定に失敗しました",
			})
			return
		}

		// Set expire date
		expire_day := time.Now().Add(time.Hour * 24 * 7)

		// update database
		result := database.Sqldb.Model(&model.User{}).Where("id = ?", token_info.UserID).Updates(map[string]interface{}{"password": json.Password, "expire_date": expire_day})
		if result.Error != nil { // error
			log.Println(result.Error)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}
		if result.RowsAffected == 0 { // mismatch user id
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}

		// response
		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
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

		//request data のjsonを変換
		var json jsonmodel.AuthLogoutRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです。",
			})
			return
		}

		access_token := token_info.Token
		device_id := token_info.DeviceID

		// logoutリクエストのため、access tokenを削除する
		if err := auth.DeleteToken(access_token, device_id); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "failed",
				"error":  "サーバでエラーが生じました。",
			})
			return
		}

		if err := database.Sqldb.Where("access_token = ?", access_token).Delete(&model.Session{}).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "failed",
				"error":  "サーバでエラーが生じました。",
			})
			return
		}

		//レスポンス
		ctx.JSON(http.StatusOK, gin.H {
			"result": "success",
		})
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "Failed to create session.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result":     "success",
			"session_id": session_id,
		})
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
			ctx.JSON(http.StatusNotAcceptable, gin.H {
				"result": "failed",
				"error":  "Invalid refresh token",
			})
			return
		}
				
		new_token, new_refresh, err := auth.SetToken(token_info.UserID, token_info.DeviceID, token_info.PublicKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "failed",
				"error": "something went wrong",
			})
			return
		}

		session_id, err := auth.CreateSession(device_id, new_token, user_id)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "Failed to create session.",
			})
			return
		}
			
		ctx.JSON(http.StatusContinue, gin.H {
			"result": "refreshed",
			"token": new_token,
			"refresh_token": new_refresh,
			"session_id": session_id,
		})
	}
}