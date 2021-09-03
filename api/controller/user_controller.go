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

		//登録しようとするユーザが既にいないか確認 (name)
		var c_name int64
		if err := database.Sqldb.Model(&model.User{}).Where("name = ?", json.Name).Select("id").Count(&c_name).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}
		if c_name > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "このユーザ名は既に使われています。",
			})
			return
		}

		//登録しようとするユーザが既にいないか確認 (email)
		var c_email int64
		if err := database.Sqldb.Model(&model.User{}).Where("email = ?", json.Email).Select("id").Count(&c_email).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}
		if c_email > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "このメールアドレスは既に使われています。",
			})
			return
		}

		//ユーザ新規作成。保存
		if err := database.Sqldb.Create(&model.User{Name: json.Name, Email: json.Email, Password: json.Password}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//新規作成したユーザのIDを取得
		var u model.User
		if err := database.Sqldb.Where("name = ? AND email = ?", json.Name, json.Email).Select("id").Take(&u).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//アクセストークンを生成
		token, err := auth.SetToken(u.ID, json.DeviceID)
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
