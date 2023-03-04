package database

import (
	model2 "github.com/Project-Birdol/birdol-server/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

/*
	TODO: Implement DB Connection Initialization
		  To use connection pool
*/
// データベースのマイグレーション -> sql.go
// DB接続はCLoseせずオブジェクトを保持 -> sql.go

func InitializeDB() *gorm.DB {
	db := getGormInstance()
	MigrateDB(db)
	return db
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&model2.User{})
	db.AutoMigrate(&model2.AccessToken{})
	db.AutoMigrate(&model2.Session{})
	db.AutoMigrate(&model2.StoryProgress{})
	db.AutoMigrate(&model2.CharacterProgress{})
	db.AutoMigrate(&model2.Teacher{})
	db.AutoMigrate(&model2.CompletedProgress{})
	db.AutoMigrate(&model2.ValidClient{})
}

func getGormInstance() *gorm.DB {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	DBADRESS := os.Getenv("DB_ADDRESS")
	PROTOCOL := "tcp(" + DBADRESS + ")"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&tls=true"
	count := 0
	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
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
	return db
}

func InitilaizeTestingDB() *gorm.DB {
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
	MigrateDB(db)
	return db
}
