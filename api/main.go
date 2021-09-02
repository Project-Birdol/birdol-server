package main

import (
	"log"
	"os"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/server"
	"github.com/gin-gonic/gin"
)

func main() {
	/*
		TODO: main.go内の処理を他ファイルに分離し整理
	*/
	database.StartDB()
	// アクセストークンの定期的な削除をする
	// auth.StartDeleteExpiredTokens()

	mode := os.Getenv("MODE")

	// ルーティング設定 -> server/server.go or server/router.go とかが多い
	var router *gin.Engine
	version := os.Getenv("API_VERSION")
	if version == "v1"{
		router = server.GetRouterV1()
	}else{
		router = server.GetRouterV2()
	}
	
	PORT := ":80"
	if mode == "production" {
		log.Println("Running in Production mode.")
	} else {
		log.Println("Running in Development mode.")
	}
	if os.Getenv("PORT") != "" {
		PORT = ":" + os.Getenv("PORT")
	}
	router.Run(PORT)
}
