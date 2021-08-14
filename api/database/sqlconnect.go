package database

import (
	"fmt"
  	"time"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)



func SqlConnect() (database *gorm.DB) {
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"
  
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	count := 0
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{})
	if err != nil {
	  for {
		fmt.Print(err) 
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