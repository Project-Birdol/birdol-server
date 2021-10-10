package model

type CharacterProgress struct
{
	Model 
    StoryProgressId int `json:"-"`
	//メインキャラクター
	MainCharacterId int	`json:"MainCharacterId"`
	Name string `json:"Name"`
	Visual float32	`json:"Visual"`
	Vocal float32	`json:"Vocal"`
	Dance float32	`json:"Dance"`
	ActiveSkillLevel uint	`json:"ActiveSkillLevel"`
	//サポートキャラクター
	SupportCharacterId uint `json:"SupportCharacterId"`
	PassiveSkillLevel uint	`json:"PassiveSkillLevel"`
	TriggeredSubStory bool `json:"TriggeredSubStory"`
}

