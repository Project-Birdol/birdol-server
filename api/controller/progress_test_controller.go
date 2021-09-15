package controller

import (
	"log"
	"net/http"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetCurrentProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid") 
		var story model.StoryProgress
		if err := database.Sqldb.Debug().Where("user_id = ? && completed = ?", userid, false).Preload("CharacterProgresses").Preload("CharacterProgresses.MainCharacter").Preload("CharacterProgresses.SupportCharacter").Last(&story).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		resonse := new(jsonmodel.StoryResponse)
		resonse.Result = "success"
		resonse.Story = story
		ctx.JSON(http.StatusOK, resonse)
	}
}

func UpdateStory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		
	}
}

func UpdateCharacters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		
	}
}


func GetGallaryInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//GORMにDISTINCTが無いみたいですね...
		/*
		DISTINCTを使う場合:SELECT DISTINCT character_id FROM birdoldb.character_progresses WHERE story_progress_id in (select id as story_progress_id FROM birdoldb.story_progresses where user_id = 1 and completed = true);
		GROUP BYを使う場合:SELECT character_id FROM birdoldb.character_progresses WHERE story_progress_id in (select id as story_progress_id FROM birdoldb.story_progresses where user_id = 1 and completed = true) GROUP BY character_id;
		*/
		userid := ctx.Param("userid")
		var ids []jsonmodel.GallaryResponse
		sub := database.Sqldb.Model(&model.StoryProgress{}).Select("ID").Where("user_id = ? && completed = ?", userid, true)
		if err := database.Sqldb.Model(&model.CharacterProgress{}).Select("character_id").Where("story_progress_id IN (?)", sub).Group("character_id").Find(&ids).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "該当する進捗が見つかりません",
			})
			return
		}
		log.Println(ids)
		ctx.JSON(http.StatusOK, ids)
		
	}
}

func GetCompletedCharacters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		/*
		SELECT * FROM birdoldb.character_progresses WHERE story_progress_id in (select id as story_progress_id FROM birdoldb.story_progresses where user_id = 1 and completed = true);
		*/
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
		log.Println(characters)
		ctx.JSON(http.StatusOK, gin.H {
			"result": "success",
		})
	}
}


func NewProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Param("userid") 
		var json jsonmodel.ProgressRequest
		if err := ctx.ShouldBindJSON(&json); err != nil && json.Characters != nil {
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
		characters := []model.CharacterProgress{}
		for _ , character := range json.Characters{
			pair := model.CharacterProgress{
				MainCharacter: model.MainCharacter{
					CharacterId: character.CharacterId,
					Visual:  character.Visual,
					Vocal:  character.Vocal,
					Dance:  character.Dance,
					ActiveSkillLevel:  character.ActiveSkillLevel,
					ActiveSkillType:  character.ActiveSkillType,
					ActiveSkillScore:  character.ActiveSkillScore,
				},
				SupportCharacter: model.SupportCharacter{
					CharacterId:  character.SupportCharacterId,
					PassiveSkillLevel:  character.PassiveSkillLevel,
					PassiveSkillType:  character.PassiveSkillType,
					PassiveSkillScore:  character.PassiveSkillScore,
				},
			}
			characters = append(characters,pair)
		}
		story := model.StoryProgress{
			UserId: uint(u64),
			CharacterProgresses: characters,
		}
		if err := database.Sqldb.Create(&story).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"result": "failed",
				"error":  "データの保存に失敗しました",
			})
			return
		}
		log.Println(json)
		ctx.JSON(http.StatusOK, gin.H {
			"result": "success",
			"session_id": "ok",
		})
	}
}
