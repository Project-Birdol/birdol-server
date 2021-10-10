package model

type StoryProgress struct
{
	Model
	UserId uint	`json:"-"`
	MainStoryId string	`json:"main_story_id"`
	Completed bool	`json:"-"`
    CharacterProgresses []CharacterProgress `gorm:"foreignKey:StoryProgressId" json:"-"`
	Teachers []Teacher `gorm:"foreignKey:StoryProgressId" json:"-"`
}
