package model

type CharacterProgress struct
{
	Model 
    StoryProgressId int `json:"-"`
	MainCharacterId int	`json:"-"`
	SupportCharacterId int	`json:"-"`
	MainCharacter MainCharacter	`json:"main_character"`
	SupportCharacter SupportCharacter	`json:"support_character"`
}
