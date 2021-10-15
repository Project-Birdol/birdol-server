package jsonmodel

import "github.com/MISW/birdol-server/database/model"

type GalleryChild struct {
	MainCharacterId int `json:"id" binding:"required"`
}

type GalleryResponse struct {
	Birdols []GalleryChild `json:"birdols"` 
}

type StoryResponse struct{
	ID uint `gorm:"primarykey" json:"id"` 
	MainStoryId string	`json:"main_story_id"`
	LessonCount uint	`json:"lesson_count"`
}

type HallOfFameResponse struct{
	Characters []model.CompletedProgress `json:"characters"` 
}

type CharacterResponse struct{
	CharacterProgresses []model.CharacterProgress `json:"character_progresses"`
	Teachers []model.Teacher `json:"teachers"`
}

type CharacterProgressRequest struct{
	CharacterProgresses []model.CharacterProgress `json:"character_progresses"`
	Teachers []model.CompletedProgress `json:"teachers"`
}

type CreateResponse struct {
	ProgressId uint	`json:"progress_id"`
	Characters []CreateCharacterChild `json:"characters"`
	Teachers []CreateTeacherChild `json:"teachers"`  
}

type CreateCharacterChild struct{
	ChracterId uint `json:"character_id"`
}

type CreateTeacherChild struct{
	TeacherId uint `json:"teacher_id"`
}