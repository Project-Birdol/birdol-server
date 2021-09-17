package model

type CharacterProgress struct
{
	Model 
    StoryProgressId int `json:"-"`
	//メインキャラクター
	MainCharacterId int	`json:"main_character_id"`
	Name string `json:"name"`
	Visual float32	`json:"visual"`
	Vocal float32	`json:"vocal"`
	Dance float32	`json:"dance"`
	ActiveSkillLevel uint	`json:"active_skill_level"`
	ActiveSkillType uint	`json:"active_skill_type"`
	ActiveSkillScore float32	`json:"active_skill_score"`
	//サポートキャラクター
	SupportCharacterId uint `json:"support_character_id"`
	PassiveSkillLevel uint	`json:"passive_skill_level"`
	PassiveSkillType uint	`json:"passive_skill_type"`
	PassiveSkillScore float32	`json:"passive_skill_score"`
}

