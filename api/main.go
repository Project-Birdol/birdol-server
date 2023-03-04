package main

import (
	"github.com/Project-Birdol/birdol-server/database"
	"github.com/Project-Birdol/birdol-server/server"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	db := database.InitializeDB()
	// アクセストークンの定期的な削除をする
	// session.StartDeleteExpiredTokens()

	mode := os.Getenv("GIN_MODE")

	// ルーティング設定
	var router *gin.Engine
	version := os.Getenv("API_VERSION")
	if version == "v1" {
		router = server.GetRouterV1(db)
	} else {
		router = server.GetRouterV2(db)
	}
	if mode == "production" {
		log.Println("Running in Production mode.")
	} else {
		log.Println("Running in Development mode.")
	}
	PORT := ":" + os.Getenv("PORT")
	router.Run(PORT)
}
