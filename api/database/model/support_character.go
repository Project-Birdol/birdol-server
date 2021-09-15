package model

type SupportCharacter struct
{
	Model 
	//サポートキャラクター
	CharacterId uint `json:"character_id"`
	PassiveSkillLevel uint	`json:"passive_skill_level"`
	PassiveSkillType uint	`json:"passive_skill_type"`
	PassiveSkillScore float32	`json:"passive_skill_score"`
}
