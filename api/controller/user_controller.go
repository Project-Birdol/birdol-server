package controller

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/MISW/birdol-server/auth"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

// 新規ユーザ登録
func HandleSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// map to struct
		var json jsonmodel.SignupUserRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}

		// デフォルトのaccount_idを生成
		i := 0
		account_id := generateRandomString(64)
		for {
			// 生成失敗
			i++
			if i > 100 {
				log.Printf("failed to create account id.")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": "failed",
					"error":  "アカウント作成に失敗しました",
				})
				return
			}

			// 重複チェック
			var c_account_id int64
			if err := database.Sqldb.Model(&model.User{}).Where("account_id = ?", account_id).Select("id").Count(&c_account_id).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": "failed",
					"error":  "アカウント作成に失敗しました",
				})
				return
			}
			if c_account_id == 0 {
				break
			}
			account_id = generateRandomString(64)
		}

		// ユーザ新規作成。保存
		new_user := model.User{Name: json.Name, AccountID: account_id, LinkPassword: model.LinkPassword{ExpireDate: time.Now()}}
		if err := database.Sqldb.Create(&new_user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		// 新規作成したユーザのIDを取得
		// log.Printf("[TEST] USER ID: %d\n", u)

		// アクセストークンを生成
		token, refresh_token, err := auth.SetToken(new_user.ID, json.DeviceID, json.PublicKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "ユーザの新規作成に失敗しました",
			})
			return
		}

		// Successful
		ctx.JSON(http.StatusOK, gin.H{
			"result":        "success",
			"user_id":       new_user.ID,
			"access_token":  token,
			"refresh_token": refresh_token,
			"account_id":    account_id,
		})
	}
}

// LinkAccount Login: account_idとpasswordで認証後にaccess tokenを発行する。 ゲーム内でのアカウント連携
func LinkAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[Login] ")
		//request data の jsonを変換
		var json jsonmodel.DataLinkRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです。",
			})
			return
		}

		//account_idが合っているかを確認。そのaccount_idでdatabaseからデータ取得
		var u model.User
		if err := database.Sqldb.Where("account_id = ?", json.AccountID).Take(&u).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "failed",
				"error":  "データ連携に失敗しました。",
			})
			return
		}

		//passwordが合っているかHash値を比較
		if err := auth.CompareHashedString(u.LinkPassword.Password, json.Password); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": "failed",
				"error":  "データ連携に失敗しました。",
			})
			return
		}

		//access token の生成及び保存
		token, refresh_token, err := auth.SetToken(u.ID, json.DeviceID, json.PublicKey)
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
			"refresh_token": refresh_token,
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
