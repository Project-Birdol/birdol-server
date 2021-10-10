package jsonmodel

import "github.com/MISW/birdol-server/database/model"

type GalleryChild struct {
	MainCharacterId int `json:"id" binding:"required"`
}

type GalleryResponse struct {
	Result string `json:"result"`
	Birdols []GalleryChild `json:"birdols"` 
}

type StoryResponse struct{
	ID uint `gorm:"primarykey" json:"id"` 
	Result string `gorm:"-" json:"result"`
	MainStoryId string	`json:"main_story_id"`
}

type DendouResponse struct{
	Result string `json:"result"`
	Characters []model.CompletedProgress `json:"characters"` 
}

type CharacterResponse struct{
	Result string `gorm:"-" json:"result"`
	CharacterProgresses []model.CharacterProgress `json:"character_progresses"`
	Teachers []model.Teacher `json:"teachers"`
}

type CharacterProgressRequest struct{
	CharacterProgresses []model.CharacterProgress `json:"character_progresses"`
	Teachers []model.CompletedProgress `json:"teachers"`
}

type CreateResponse struct {
	Result string `json:"result"`
	ProgressId uint	`json:"progress_id"`
	Characters []CreateCharacterChild `json:"characters"`
	Teachers []CreateTeacherChild `json:"teachers"`  
}

type CreateCharacterChild struct{
	ChracterId uint `json:"chracter_id"`
}

type CreateTeacherChild struct{
	TeacherId uint `json:"teacher_id"`
}