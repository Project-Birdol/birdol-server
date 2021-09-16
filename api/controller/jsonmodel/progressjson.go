package jsonmodel

import "github.com/MISW/birdol-server/database/model"

type GallaryChild struct {
	MainCharacterId int `json:"id" binding:"required"`
}

type GallaryResponse struct {
	Result string `json:"result"`
	Birdols []GallaryChild `json:"birdols"` 
}

type StoryResponse struct{
	Result string `json:"result"`
	Story model.StoryProgress `json:"story_progress"` 
}

type DendouResponse struct{
	Result string `json:"result"`
	Characters []model.CompletedProgress `json:"characters"` 
}