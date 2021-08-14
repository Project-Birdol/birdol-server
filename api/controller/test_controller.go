package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/MISW/birdol-server/controller/jsonmodel"
)

//GET:取得　POST:情報の送信&データの新規作成 PUT:データの更新 DELETE:データの削除

func TestGet() gin.HandlerFunc {
    return func(ctx *gin.Context) {
		sqldb := database.SqlConnect()
    	var users []model.User
    	sqldb.Order("created_at asc").Find(&users)
		db, _ := sqldb.DB()
    	defer db.Close()
        ctx.JSON(http.StatusOK, gin.H{
			"users": users,
		})
    }
}

func TestPost() gin.HandlerFunc {
    return func(ctx *gin.Context) {
		sqldb := database.SqlConnect()
		var json jsonmodel.UserRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
    	result := sqldb.Create(&model.User{Name: json.Name, Email: json.Email})
    	db, _ := sqldb.DB()
		defer db.Close()
		fmt.Println(result.Error)
		if result.Error != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
			})
		}else{
			ctx.JSON(http.StatusOK, gin.H{
				"result": "success",
			})
		}
		
    }
}

func TestPut() gin.HandlerFunc {
    return func(ctx *gin.Context) {
		sqldb := database.SqlConnect()
		n := ctx.Param("id")
    	id, err := strconv.Atoi(n)
    	if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
    	}
		var json jsonmodel.UserRequest
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user model.User
    	sqldb.First(&user, id)
		user.Name = json.Name
		user.Email = json.Email
    	result := sqldb.Save(&user)
		db, _ := sqldb.DB()
    	defer db.Close()
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		}else{
			ctx.JSON(http.StatusOK, gin.H{
				"result": "success",
			})
		}
    }
}

func TestDelete() gin.HandlerFunc {
    return func(ctx *gin.Context) {
		sqldb := database.SqlConnect()
		n := ctx.Param("id")
    	id, err := strconv.Atoi(n)
    	if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
    	}
		result := sqldb.Delete(&model.User{},id)
		db, _ := sqldb.DB()
    	defer db.Close()
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
		
    }
}