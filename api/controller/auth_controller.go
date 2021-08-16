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
		//datanase connection
		sqldb := database.SqlConnect()
		db, _ := sqldb.DB()
		defer db.Close()

		//request data の jsonを変換
		var json jsonmodel.AuthLoginRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//emailとpasswordが合っているかを確認
		var u model.User
		if err := sqldb.Where("email = ? AND password = ?", json.Email, json.Password).Take(&u).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		//access token の生成及び保存
		token, err := auth.SetToken(sqldb, u.ID)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		//database connection
		sqldb := database.SqlConnect()
		db, _ := sqldb.DB()
		defer db.Close()

		//request data のjsonを変換
		var json jsonmodel.AuthLogoutRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//access token が正しいか確認
		if err := auth.CheckToken(sqldb, json.Auth.UserID, json.Auth.AccessToken); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		//logoutリクエストのため、access tokenを削除する。
		if err := auth.DeleteToken(sqldb, json.Auth.UserID); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//レスポンス
		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
	}
}
