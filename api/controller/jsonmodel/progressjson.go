package jsonmodel

import "github.com/MISW/birdol-server/database/model"

type ProgressRequest struct {
	Characters []Character `json:"characters" binding:"dive"`
	Teachers []Teacher `json:"teachers"`
}

type GallaryResponse struct {
	CharacterId int `json:"id" binding:"required"`
}

type StoryResponse struct{
	Result string `json:"result"`
	Story model.StoryProgress `json:"story_progress"` 
}

type DendouResponse struct{
	Result string `json:"result"`
	Pairs model.CharacterProgress `json:"pairs"` 
}