package controller

import (
	"net/http"

	"github.com/MISW/birdol-server/auth"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

//新規ユーザ登録
func HandleSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//database
		sqldb := database.SqlConnect()
		db, _ := sqldb.DB()
		defer db.Close()

		//request data のjsonをstructにする
		var json jsonmodel.SignupUserRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}

		//request data に含まれるパスワードをハッシュ化する
		if err := auth.HashString(&json.Password); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//登録しようとするユーザが既にいないか確認
		var c int64
		if err := sqldb.Model(&model.User{}).Where("name = ? AND email = ?", json.Name, json.Email).Select("id").Count(&c).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}
		if c > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "このアカウントは既に存在しています。",
			})
			return
		}

		//ユーザ新規作成。保存
		if err := sqldb.Create(&model.User{Name: json.Name, Email: json.Email, Password: json.Password}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//新規作成したユーザのIDを取得
		var u model.User
		if err := sqldb.Where("name = ? AND email = ?", json.Name, json.Email).Select("id").Take(&u).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//アクセストークンを生成
		token, err := auth.SetToken(sqldb, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//ユーザ新規登録成功!
		ctx.JSON(http.StatusOK, gin.H{
			"result":       "success",
			"user_id":      u.ID,
			"access_token": token,
		})
	}
}
