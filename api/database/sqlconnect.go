package database

import (
	"fmt"
  	"time"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"os"
)



func SqlConnect() (database *gorm.DB) {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	DBADRESS := os.Getenv("DB_ADDRESS")
	PROTOCOL := "tcp(" + DBADRESS + ":3306)"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true"
	count := 0
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{})
	if err != nil {
	  for {
		fmt.Print(err)
		fmt.Print(CONNECT) 
		if err == nil {
		  fmt.Println("")
		  break
		}
		fmt.Print(".")
		time.Sleep(time.Second)
		count++
		if count > 180 {
		  fmt.Println("DB Connection Error")
		  panic(err)
		}
		db, err = gorm.Open(mysql.Open(CONNECT), &gorm.Config{})
	  }
	}
	fmt.Println("DB Connected")
  
	return db
}