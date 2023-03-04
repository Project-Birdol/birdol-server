package controller

import (
	"encoding/json"
	"errors"
	"github.com/Project-Birdol/birdol-server/controller/jsonmodel"
	"github.com/Project-Birdol/birdol-server/model"
	res_util "github.com/Project-Birdol/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type ProgressController struct {
	DB *gorm.DB
}

func (pc *ProgressController) GetGalleryInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID
		var ids []jsonmodel.GalleryChild
		if err := pc.DB.Model(&model.CompletedProgress{}).Select("main_character_id").Where("user_id = ?", userid).Group("main_character_id").Order("main_character_id").Find(&ids).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrDataNotFound)
			return
		}

		response := jsonmodel.GalleryResponse{
			Birdols: ids,
		}
		res_util.SetStructResponse(ctx, http.StatusOK, res_util.ResultOK, response)
	}
}

func (pc *ProgressController) GetCompletedCharacters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID
		var characters []model.CompletedProgress
		if err := pc.DB.Model(&model.CompletedProgress{}).Where("user_id = ?", userid).Find(&characters).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrDataNotFound)
			return
		}

		response := jsonmodel.HallOfFameResponse{
			Characters: characters,
		}
		log.Println(response)
		res_util.SetStructResponse(ctx, http.StatusOK, res_util.ResultOK, response)
	}
}

func (pc *ProgressController) GetCurrentCharacters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID

		var story model.StoryProgress
		if err := pc.DB.Where("user_id = ? && completed = ?", userid, false).Preload("CharacterProgresses").Preload("Teachers").Preload("Teachers.Character").Last(&story).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrDataNotFound)
			return
		}

		response := jsonmodel.CharacterResponse{
			CharacterProgresses: story.CharacterProgresses,
			Teachers:            story.Teachers,
		}
		res_util.SetStructResponse(ctx, http.StatusOK, res_util.ResultOK, response)
	}
}

func (pc *ProgressController) GetCurrentStory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID
		log.Println(userid)
		var response jsonmodel.StoryResponse
		if err := pc.DB.Model(&model.StoryProgress{}).Where("user_id = ? && completed = ?", userid, false).Last(&response).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrDataNotFound)
			return
		}
		res_util.SetStructResponse(ctx, http.StatusOK, res_util.ResultOK, response)
	}
}

func (pc *ProgressController) FinishProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID
		var story model.StoryProgress
		if err := pc.DB.Where("user_id = ? && completed = ?", userid, false).Last(&story).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrDataNotFound)
			return
		}
		story.Completed = true
		if err := pc.DB.Save(&story).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusInternalServerError, res_util.ErrFailDataStore)
			return
		}
		var characters []model.CompletedProgress
		if err := pc.DB.Model(&model.CharacterProgress{}).Select("main_character_id", "name", "visual", "vocal", "dance", "active_skill_level", "support_character_id", "passive_skill_level").Where("story_progress_id = ?", story.ID).Find(&characters).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusInternalServerError, res_util.ErrFailDataStore)
			return
		}
		for i := 0; i < 5; i++ {
			characters[i].UserId = userid
		}
		if err := pc.DB.Model(&model.CompletedProgress{}).Create(&characters).Error; err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusInternalServerError, res_util.ErrFailDataStore)
			return
		}
		res_util.SetNormalResponse(ctx, http.StatusOK, res_util.ResultOK)
	}
}

func (pc *ProgressController) UpdateCharacters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 育成進捗サブシナリオのステータスの更新
		var request jsonmodel.CharacterProgressRequest
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &request); err != nil {
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrFailParseJSON)
			return
		}
		for _, v := range request.CharacterProgresses {
			if err := pc.DB.Model(&model.CharacterProgress{}).Where("id = ?", v.ID).Updates(&v).Error; err != nil {
				res_util.SetErrorResponse(ctx, http.StatusInternalServerError, res_util.ErrFailDataStore)
				return
			}
		}
		res_util.SetNormalResponse(ctx, http.StatusOK, res_util.ResultOK)
	}
}

func (pc *ProgressController) UpdateMainStory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID
		var request model.StoryProgress
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &request); err != nil {
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrFailParseJSON)
			return
		}
		if err := pc.DB.Model(&model.StoryProgress{}).Where("user_id = ? && completed = ?", userid, false).Updates(&request).Error; err != nil {
			res_util.SetErrorResponse(ctx, http.StatusInternalServerError, res_util.ErrFailDataStore)
			return
		}
		res_util.SetNormalResponse(ctx, http.StatusOK, res_util.ResultOK)
	}
}

func (pc *ProgressController) CreateProgress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, _ := ctx.Get("access_token")
		userid := accessToken.(model.AccessToken).UserID
		var request jsonmodel.CharacterProgressRequest
		body_byte_interface, _ := ctx.Get("body_rawbyte")
		body_rawbyte := body_byte_interface.([]byte)
		if err := json.Unmarshal(body_rawbyte, &request); err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusBadRequest, res_util.ErrFailParseJSON)
			return
		}
		log.Println(request.Teachers)

		var err error
		var current []model.StoryProgress
		pc.DB.Where("user_id = ? && completed = ?", userid, false).Find(&current)
		proglength := len(request.CharacterProgresses)
		teachlength := len(request.Teachers)
		story := model.StoryProgress{}
		story.UserId = userid
		story.CharacterProgresses = request.CharacterProgresses
		story.Teachers = []model.Teacher{}
		for _, newteacher := range request.Teachers {
			story.Teachers = append(story.Teachers, model.Teacher{
				CharacterId: newteacher.ID,
				Character:   newteacher,
			})
		}
		if len(current) != 0 || (proglength > 0 && proglength != 5) || (teachlength > 0 && teachlength != 1) {
			err = errors.New("invalid request")
		} else {
			//進捗の新規作成
			err = pc.DB.Create(&story).Error
		}

		if err != nil {
			log.Println(err)
			res_util.SetErrorResponse(ctx, http.StatusInternalServerError, res_util.ErrFailDataStore)
			return
		}
		characters := []jsonmodel.CreateCharacterChild{}
		log.Println(story)
		for _, character := range story.CharacterProgresses {
			characters = append(characters, jsonmodel.CreateCharacterChild{
				ChracterId: character.ID,
			})
		}
		teachers := []jsonmodel.CreateTeacherChild{}
		for _, teacher := range story.Teachers {
			teachers = append(teachers, jsonmodel.CreateTeacherChild{
				TeacherId: teacher.ID,
			})
		}
		response := jsonmodel.CreateResponse{
			ProgressId: story.ID,
			Characters: characters,
			Teachers:   teachers,
		}
		res_util.SetStructResponse(ctx, http.StatusOK, res_util.ResultOK, response)
	}
}
