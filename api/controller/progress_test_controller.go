package controller

import (
	"log"
	"net/http"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"errors"
	"strconv"
)

func GetCurrentProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid") 
		var story model.StoryProgress
		if err := database.Sqldb.Where("user_id = ? && completed = ?", userid, false).Preload("CharacterProgresses").Preload("CharacterProgresses.MainCharacter").Preload("CharacterProgresses.SupportCharacter").Preload("Teachers").Preload("Teachers.Character").Last(&story).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		response := new(jsonmodel.StoryResponse)
		response.Result = "success"
		response.Story = story
		ctx.JSON(http.StatusOK, response)
	}
}

func GetGallaryInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		var ids []jsonmodel.GallaryChild
		sub1 := database.Sqldb.Model(&model.StoryProgress{}).Select("ID").Where("user_id = ? && completed = ?", userid, true)
		sub2 := database.Sqldb.Model(&model.CharacterProgress{}).Select("ID").Where("story_progress_id IN (?)", sub1)
		if err := database.Sqldb.Model(&model.MainCharacter{}).Select("character_id").Where("ID IN (?)", sub2).Group("character_id").Order("character_id").Find(&ids).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		response := new(jsonmodel.GallaryResponse)
		response.Result = "success"
		response.Birdols = ids
		ctx.JSON(http.StatusOK, response)
		
	}
}

func GetCompletedCharacters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		var characters []model.CharacterProgress
		sub := database.Sqldb.Model(&model.StoryProgress{}).Select("ID").Where("user_id = ? && completed = ?", userid, true)
		if err := database.Sqldb.Where("story_progress_id IN (?)", sub).Find(&characters).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		response := new(jsonmodel.DendouResponse)
		response.Result = "success"
		response.Pairs = characters
		ctx.JSON(http.StatusOK, response)
	}
}


func CreateOrUpdateProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid") 
		var story model.StoryProgress
		if err := ctx.ShouldBindJSON(&story); err != nil || story.CharacterProgresses == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです",
			})
			return
		}
		u64, uerr := strconv.ParseUint(userid, 10, 32)
		if uerr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データの保存に失敗しました",
			})
			return
		}
		story.UserId = uint(u64)
		var err error
		if story.ID != 0{
			var newstory model.StoryProgress
			database.Sqldb.Where("user_id = ? && completed = ?", userid, false).Last(&newstory)
			if newstory.ID != story.ID{
				err = errors.New("Invalid request")
			}else{
				//進捗の更新
				err = database.Sqldb.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&story).Error
			}
		}else{
			//進捗の新規作成
			//先生のチェック
			err = database.Sqldb.Create(&story).Error;		
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データの保存に失敗しました",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H {
			"result": "success",
			"session_id": "ok",
		})
	}
}
