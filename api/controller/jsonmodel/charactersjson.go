package jsonmodel

type Character struct {
	CharacterId      int   `json:"character_id" binding:"required"`
	Visual      float32   `json:"visual"`
	Vocal      float32   `json:"vocal"`
	Dance      float32   `json:"dance"`
	ActiveSkillLevel      uint   `json:"active_skill_level"`
	ActiveSkillType      uint   `json:"active_skill_type"`
	ActiveSkillScore      float32   `json:"active_skill_score"`
	SupportCharacterId      uint   `json:"support_character_id"`
	PassiveSkillLevel      uint   `json:"passive_skill_level"`
	PassiveSkillType      uint   `json:"passive_skill_type"`
	PassiveSkillScore      float32   `json:"passive_skill_score"`
}



type Teacher struct {
	SupportCharacterId      uint   `json:"support_character_id" binding:"required"`
	PassiveSkillLevel      uint   `json:"passive_skill_level" binding:"required"`
	PassiveSkilltype      uint   `json:"passive_skill_type" binding:"required"`
	PassiveSkillScore      float32   `json:"passive_skill_score" binding:"required"`
}

