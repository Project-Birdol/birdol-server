package main

import (
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
)



func main() {
	sqldb := database.SqlConnect()
 	sqldb.AutoMigrate(&model.User{})
	db, err := sqldb.DB()
  	defer db.Close()
	if err != nil {
		fmt.Println("Database Error:")
	}
	router := gin.Default()
	//　ルーティングはここで設定する
	router.GET("/api/v1/test",controller.TestGet())
	router.POST("/api/v1/test",controller.TestPost())
	router.PUT("/api/v1/test/:id",controller.TestPut())
	router.DELETE("/api/v1/test/:id",controller.TestDelete())
	mode := os.Getenv("MODE")
	PORT := ":80"
	if mode=="production"{
		fmt.Println("Running in Production mode.")
	}else{
		fmt.Println("Running in Development mode.")
	}
	if os.Getenv("PORT")!=""{
		PORT = ":"+os.Getenv("PORT")
	}
	router.Run(PORT)
}



