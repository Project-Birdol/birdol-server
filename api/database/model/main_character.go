package model

type MainCharacter struct
{
	Model 
	//メインキャラクター
    CharacterId int	`json:"character_id"`
	Name string `json:"name"`
	Visual float32	`json:"visual"`
	Vocal float32	`json:"vocal"`
	Dance float32	`json:"dance"`
	ActiveSkillLevel uint	`json:"active_skill_level"`
	ActiveSkillType uint	`json:"active_skill_type"`
	ActiveSkillScore float32	`json:"active_skill_score"`
}
