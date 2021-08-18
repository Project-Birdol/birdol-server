package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

func main() {
	// データベースのマイグレーション
	sqldb := database.SqlConnect()
	sqldb.AutoMigrate(&model.User{})
	sqldb.AutoMigrate(&model.AccessToken{})
	sqldb.AutoMigrate(&model.Session{})
	db, err := sqldb.DB()
	if err != nil {
		log.Fatal("Database error: ", err)
	}
	defer db.Close()

	// アクセストークンの定期的な削除をする
	// auth.StartDeleteExpiredTokens()

	mode := os.Getenv("MODE")

	router := gin.Default()
	/*	ルーティングはここで設定する
		  v1: For development
		  v2: For production	*/
	v1 := router.Group("api/v1")
	{
		user := v1.Group("/user")
		{
			user.PUT("", controller.HandleSignUp())
			user.POST("", controller.HandleLogin())
			user.DELETE("", controller.HandleLogout())
		}
	}
	
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
