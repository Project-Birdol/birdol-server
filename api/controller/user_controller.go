package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/MISW/birdol-server/auth"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/MISW/birdol-server/utils/response"
	"github.com/MISW/birdol-server/utils/random"
	"github.com/gin-gonic/gin"
)

// 新規ユーザ登録
func HandleSignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[SignUp] ")

		content_type := ctx.GetHeader("Content-Type")
		if content_type != gin.MIMEJSON {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidType)
			ctx.Abort()
			return
		}

		// parse json
		var request_body jsonmodel.SignupUserRequest
		if err := ctx.ShouldBindJSON(&request_body); err != nil {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrFailParseJSON)
			return
		}

		// デフォルトのaccount_idを生成
		account_id, err := random.GenerateRandomString(9)
		if err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailAccountCreation)
			return
		}
		
		// ユーザ新規作成。保存
		new_user := model.User{Name: request_body.Name, AccountID: account_id, LinkPassword: model.LinkPassword{ExpireDate: time.Now()}}
		if err := database.Sqldb.Create(&new_user).Error; err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailAccountCreation)
			return
		}

		// アクセストークンを生成
		token, refresh_token, err := auth.SetToken(new_user.ID, request_body.DeviceID, request_body.PublicKey)
		if err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailAccountCreation)
			return
		}

		// Successful
		property := gin.H {
			"user_id":       new_user.ID,
			"access_token":  token,
			"refresh_token": refresh_token,
			"account_id":    account_id,
		}
		response.SetNormalResponse(ctx, http.StatusOK, response.ResultOK, property)
	}
}

// LinkAccount Login: account_idとpasswordで認証後にaccess tokenを発行する。 ゲーム内でのアカウント連携
func LinkAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[DataLink] ")

		content_type := ctx.GetHeader("Content-Type")
		if content_type != gin.MIMEJSON {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidType)
			ctx.Abort()
			return
		}

		// parse json
		var request_body jsonmodel.DataLinkRequest
		if err := ctx.ShouldBindJSON(&request_body); err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrFailParseJSON)
			return
		}

		// account_idが合っているかを確認。そのaccount_idでdatabaseからデータ取得
		var u model.User
		if err := database.Sqldb.Where("account_id = ?", request_body.AccountID).Take(&u).Error; err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusNotFound, response.ErrInvalidAccount)
			return
		}

		// expire check
		now := time.Now()
		if now.After(u.LinkPassword.ExpireDate) {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrPasswordExpire)
			return
		}

		// passwordが合っているかHash値を比較
		if err := auth.CompareHashedString(u.LinkPassword.Password, request_body.Password); err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidPassword)
			return
		}

		// disable used password
		if err := database.Sqldb.Model(&model.User{}).Where("id = ?", u.ID).Update("expire_date", time.Now()).Error; err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailDataLink)
			return
		}

		// access token の生成及び保存
		token, refresh_token, err := auth.SetToken(u.ID, request_body.DeviceID, request_body.PublicKey)
		if err != nil {
			log.Println(err)
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailDataLink)
			return
		}

		// response
		property := gin.H {
			"user_id":      u.ID,
			"access_token": token,
			"refresh_token": refresh_token,
		}
		response.SetNormalResponse(ctx, http.StatusOK, response.ResultOK, property)
	}
}
