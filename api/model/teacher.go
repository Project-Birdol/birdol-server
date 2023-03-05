package model

type Teacher struct {
	Model
	//育成に参加する殿堂入りバードル
	StoryProgressId int               `json:"-"`
	CharacterId     uint              `json:"CharacterId"`
	Character       CompletedProgress `json:"character"`
}
