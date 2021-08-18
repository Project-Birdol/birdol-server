package main

import (
	"fmt"
	"os"

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

	router.PUT("/api/v1/user", controller.HandleSignUp())
	router.POST("/api/v1/user", controller.HandleLogin())
	router.DELETE("/api/v1/auth", controller.HandleLogout())
	mode := os.Getenv("MODE")
	PORT := ":80"
	if mode == "production" {
		fmt.Println("Running in Production mode.")
	} else {
		fmt.Println("Running in Development mode.")
	}
	if os.Getenv("PORT") != "" {
		PORT = ":" + os.Getenv("PORT")
	}
	router.Run(PORT)
}
