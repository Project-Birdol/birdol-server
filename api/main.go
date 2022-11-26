package main

import (
	"log"
	"os"
	"github.com/Project-Birdol/birdol-server/database"
	"github.com/Project-Birdol/birdol-server/server"
	"github.com/gin-gonic/gin"
)

func main() {
	database.StartDB()
	// アクセストークンの定期的な削除をする
	// auth.StartDeleteExpiredTokens()

	mode := os.Getenv("GIN_MODE")

	// ルーティング設定
	var router *gin.Engine
	version := os.Getenv("API_VERSION")
	if version == "v1"{
		router = server.GetRouterV1()
	}else{
		router = server.GetRouterV2()
	}
	if mode == "production" {
		log.Println("Running in Production mode.")
	} else {
		log.Println("Running in Development mode.")
	}
	PORT := ":" + os.Getenv("PORT")
	router.Run(PORT)
}
