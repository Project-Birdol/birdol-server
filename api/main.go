package main

import (
	"fmt"

	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

func main() {
	//データベースのマイグレーション
	sqldb := database.SqlConnect()
	sqldb.AutoMigrate(&model.User{})
	sqldb.AutoMigrate(&model.AccessToken{})
	db, err := sqldb.DB()
	defer db.Close()
	if err != nil {
		fmt.Println("Database Error:")
	}

	//アクセストークンの定期的な削除をする
	//auth.StartDeleteExpiredTokens()

	router := gin.Default()
	//　ルーティングはここで設定する
	router.GET("/api/v1/test", controller.TestGet())
	router.POST("/api/v1/test", controller.TestPost())
	router.PUT("/api/v1/test/:id", controller.TestPut())
	router.DELETE("/api/v1/test/:id", controller.TestDelete())

	router.POST("/api/v1/auth", controller.HandleLogin())
	router.DELETE("/api/v1/auth", controller.HandleLogout())

	// For Development
	router.Run(":80")
}
