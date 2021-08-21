package controller

import (
	"log"
	"net/http"

	"github.com/MISW/birdol-server/auth"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

//HandleLogin Login: emailとpasswordで認証後にaccess tokenを発行する
//e.g. REQUEST: curl -X POST --data '{"email":"test@test","password":"test"}'  -H "Content-Type: application/json" http://localhost:80/api/v1/auth
//e.g. RESPONSE: {"access_token":"WXgRCCTFhR8nY1MEKv5s1nXrRfCPUVza","result":"success","user_id":11}
func HandleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[HandleLogin]")
		//datanase connection
		sqldb := database.SqlConnect()
		db, _ := sqldb.DB()
		defer db.Close()

		//request data の jsonを変換
		var json jsonmodel.AuthLoginRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです。",
			})
			return
		}

		//emailが合っているかを確認。そのemailでdatabaseからデータ取得
		var u model.User
		if err := sqldb.Where("email = ?", json.Email).Take(&u).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "failed",
				"error":  "ログインに失敗しました。", //またはそのemailのユーザが存在しないことを示す。
			})
			return
		}

		//passwordが合っているかHash値を比較
		if err := auth.CompareHashedString(u.Password, json.Password); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "failed",
				"error":  "ログインに失敗しました。",
			})
			return
		}

		//access token の生成及び保存
		token, err := auth.SetToken(sqldb, u.ID)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "サーバでエラーが生じました。",
			})
			return
		}

		//response
		ctx.JSON(http.StatusOK, gin.H{
			"result":       "success",
			"user_id":      u.ID,
			"access_token": token,
		})
	}
}

//HandleLogout Logout: user_idとaccess_tokenで認証した後にaccess_tokenを削除する。
//e.g. REQUEST: curl -X DELETE --data '{"auth":{"user_id":11,"access_token":"USACD7zX3IgiYnp4u9bSNtPOr92Pyj9N"}}' -H "Content-Type: application/json" http://localhost:80/api/v1/auth
//e.g. RESPONSE: {"result":"success"}
func HandleLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[HandleLogout]")
		//database connection
		sqldb := database.SqlConnect()
		db, _ := sqldb.DB()
		defer db.Close()

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

		//access token が正しいか確認
		if err := auth.CheckToken(sqldb, json.Auth.UserID, json.Auth.AccessToken); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "failed",
				"error":  "認証に失敗しました。",
			})
			return
		}

		//logoutリクエストのため、access tokenを削除する。
		if err := auth.DeleteToken(sqldb, json.Auth.UserID); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "サーバでエラーが生じました。",
			})
			return
		}

		//レスポンス
		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
	}
}

// TODO: Create Authorization Handler (For Login by using Token)
/* 
	認証の構成としては，一番外枠にメアパス認証があり，一度メアパスを通したら有効期間内はトークン認証
	メアパス認証時にはデバイスIDに固有のトークンを付与し，複数端末のログインを管理

	また，複数端末の同時ログインの管理として，セッションIDの付与により同一アカウントに対して複数セッションが立つのを阻止
	セッションIDはトークン認証時（アプリ起動時）に付与し，認証のたびに以前のものは無効化され，新規生成される
*/

func TokenAuthorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[TokenAuthorize]")
		// Establish Database Connection (廃止予定)
		db := database.SqlConnect()
		sqldb, _ := db.DB()
		defer sqldb.Close()

		// Processing request
		var request jsonmodel.Auth
		if err := ctx.ShouldBindJSON(&request); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H {
				"result": "failed",
				"error": "Invalid Request.",
			})
			return
		}

		user_id := request.UserID
		access_token := request.AccessToken
		device_id := request.DeviceID

		if err := auth.CheckToken(db, user_id, access_token); err != nil {
			log.Println(err)
			// TODO: if token expired, Regenerate token
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "failed",
				"error": "Invaild AccessToken.",
			})
			return
		}

		session_id, err := auth.CreateSession(db, device_id, access_token, user_id)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "failed",
				"error": "Failed to create session.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H {
			"result": "success",
			"session_id": session_id,
		})
	}
}
