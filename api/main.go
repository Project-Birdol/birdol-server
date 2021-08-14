package main

import (
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
	// For Development
	router.Run(":80")
}



