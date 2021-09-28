package database

import (
	"log"
  	"time"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/mysql"
	"github.com/MISW/birdol-server/database/model"
	"os"
)

/*
	TODO: Implement DB Connection Initialization
		  To use connection pool
*/
// データベースのマイグレーション -> sql.go
// DB接続はCLoseせずオブジェクトを保持 -> sql.go

var Sqldb *gorm.DB

func StartDB(){
	SqlConnect()
	MigrateDB()
}

func MigrateDB(){
	Sqldb.AutoMigrate(&model.User{})
	Sqldb.AutoMigrate(&model.AccessToken{})
	Sqldb.AutoMigrate(&model.Session{})
	Sqldb.AutoMigrate(&model.StoryProgress{})
	Sqldb.AutoMigrate(&model.CharacterProgress{})
	Sqldb.AutoMigrate(&model.Teacher{})
	Sqldb.AutoMigrate(&model.CompletedProgress{})
}

func SqlConnect() {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	DBADRESS := os.Getenv("DB_ADDRESS")
	PROTOCOL := "tcp(" + DBADRESS + ":3306)"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true"
	count := 0
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
	  for {
		log.Print(err)
		log.Print(CONNECT) 
		if err == nil {
		  log.Println("")
		  break
		}
		log.Print(".")
		time.Sleep(time.Second)
		count++
		if count > 180 {
		log.Println("DB Connection Error")
		  panic(err)
		}
		db, err = gorm.Open(mysql.Open(CONNECT), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	  }
	}
	log.Println("DB Connected")
	Sqldb = db
}

func TestingDatabase() {
	MigrateDB()

	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	DBADRESS := os.Getenv("DB_ADDRESS")
	PROTOCOL := "tcp(" + DBADRESS + ":3306)"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true"
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
	  panic(err)
	}
	log.Println("DB Connected")
	Sqldb = db
}