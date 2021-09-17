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
		if err := database.Sqldb.Where("user_id = ? && completed = ?", userid,false).Preload("CharacterProgresses").Preload("Teachers").Preload("Teachers.Character").Last(&story).Error; err != nil {
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

func GetGalleryInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		var ids []jsonmodel.GallaryChild
		if err := database.Sqldb.Model(&model.CompletedProgress{}).Select("main_character_id").Where("user_id = ?", userid).Group("main_character_id").Order("main_character_id").Find(&ids).Error; err != nil {
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
		var characters []model.CompletedProgress
		if err := database.Sqldb.Model(&model.CompletedProgress{}).Where("user_id = ?", userid).Find(&characters).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		response := new(jsonmodel.DendouResponse)
		response.Result = "success"
		response.Characters = characters
		ctx.JSON(http.StatusOK, response)
	}
}

func FinishProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		var user model.User
		if err := database.Sqldb.Where("id = ?", userid).Take(&user).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当するユーザーが見つかりません",
			})
			return
		}
		var story model.StoryProgress
		if err := database.Sqldb.Where("user_id = ? && completed = ?", userid,false).Last(&story).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		story.Completed = true
		if err := database.Sqldb.Save(&story).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データの保存に失敗しました",
			})
			return
		}
		var characters []model.CompletedProgress
		if err := database.Sqldb.Model(&model.CharacterProgress{}).Select("main_character_id","name","visual","vocal","dance","active_skill_level","active_skill_type","active_skill_score","support_character_id","passive_skill_level","passive_skill_type","passive_skill_score").Where("story_progress_id = ?", story.ID).Find(&characters).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データの保存に失敗しました",
			})
			return
		}
		if err := database.Sqldb.Model(&user).Association("CompletedProgresses").Append(&characters); err != nil {
			log.Println(err)
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
			var laststory model.StoryProgress
			database.Sqldb.Where("user_id = ? && completed = ?", userid, false).Last(&laststory)
			if laststory.ID != story.ID{
				err = errors.New("Invalid request")
			}else{
				//進捗の更新
				err = database.Sqldb.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&story).Error
			}
		}else{
			var stories []model.StoryProgress
			database.Sqldb.Where("user_id = ? && completed = ?", userid, false).Find(&stories)
			proglength := len(story.CharacterProgresses)
			teachlength := len(story.Teachers)
			if len(stories) != 0 || (proglength > 0 && proglength != 5) || (teachlength > 0 && teachlength != 2){
				err = errors.New("Invalid request")
			}else{
				//進捗の新規作成
				err = database.Sqldb.Create(&story).Error;	
			}
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
			"progress_id": story.ID,
		})
	}
}
