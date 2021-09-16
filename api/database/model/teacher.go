package model

type Teacher struct
{
	Model 
	//育成に参加する殿堂入りバードル
    StoryProgressId int `json:"-"`
	CharacterId int `json:"completed_character_id"`
	Character CompletedProgress `json:"character"`
}
