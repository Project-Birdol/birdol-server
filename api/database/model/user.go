package model

type User struct {
	Model
	Name                string              `gorm:"unique;not null"`
	AccountID           string              `gorm:"unique;not null"`
	Password            string              `gorm:"not null"`
	CompletedProgresses []CompletedProgress `gorm:"foreignKey:UserId" json:"completed_progresses"`
}
