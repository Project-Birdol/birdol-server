package model

type StoryProgress struct
{
	Model
	UserId uint	`json:"user_id"`
	MainStoryId uint	`json:"main_story_id"`
	Completed bool	`json:"-"`
    CharacterProgresses []CharacterProgress `gorm:"foreignKey:StoryProgressId" json:"character_progresses"`
	Teachers []Teacher `gorm:"foreignKey:StoryProgressId"`
}
