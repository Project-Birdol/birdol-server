package model

type StoryProgress struct {
	Model
	UserId              uint                `json:"-"`
	MainStoryId         string              `gorm:"default:1a" json:"main_story_id"`
	Completed           bool                `json:"-"`
	LessonCount         uint                `gorm:"default:5" json:"lesson_count"`
	CharacterProgresses []CharacterProgress `gorm:"foreignKey:StoryProgressId" json:"-"`
	Teachers            []Teacher           `gorm:"foreignKey:StoryProgressId" json:"-"`
}
