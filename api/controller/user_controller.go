package controller

import (
	"log"
	"math/rand"
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

		//デフォルトのaccount_idを生成
		i := 0
		account_id := generateRandomString(64)
		for {
			//生成失敗
			i++
			if i > 100 {
				log.Println("デフォルトのaccount_idの生成に失敗しました。")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": "failed",
					"error":  "データ連携の設定に失敗しました",
				})
				return
			}

			//重複チェック
			var c_account_id int64
			if err := database.Sqldb.Model(&model.User{}).Where("account_id = ?", account_id).Select("id").Count(&c_account_id).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": "failed",
					"error":  "データ連携の設定に失敗しました",
				})
				return
			}
			if c_account_id == 0 {
				break
			}
			account_id = generateRandomString(64)
		}

		//ユーザ新規作成。保存
		u := model.User{Name: json.Name, AccountID: account_id, Password: json.Password}
		if err := database.Sqldb.Create(&u).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		//新規作成したユーザのIDを取得
		//log.Printf("[TEST] USER ID: %d\n", u)

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
			"account_id":   account_id,
		})
	}
}

// EditAccount Loginするためのaccount_idとpasswordを設定、編集する
func EditAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		//request data の jsonを変換
		var json jsonmodel.EditAccountRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}

		//登録しようとするユーザが既にいないか確認 (account_id)
		var c_account_id int64
		if err := database.Sqldb.Model(&model.User{}).Where("account_id = ?", json.AccountID).Select("id").Count(&c_account_id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データ連携の設定に失敗しました",
			})
			return
		}
		if c_account_id > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "このidは既に使われています",
			})
			return
		}

		//request data に含まれるパスワードをハッシュ化する
		if err := auth.HashString(&json.Password); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データ連携の設定に失敗しました",
			})
			return
		}

		// update database
		result := database.Sqldb.Model(&model.User{}).Where("id = ?", userid).Updates(map[string]interface{}{"account_id": json.AccountID, "password": json.Password})
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

		//response
		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
	}
}

func generateRandomString(length int) string {
	const charas = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	clen := len(charas)
	r := make([]byte, length)
	for i := range r {
		r[i] = charas[rand.Intn(clen)]
	}
	return string(r)
}
