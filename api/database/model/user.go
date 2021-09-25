package model

type User struct {
	Model
	Name                string              `gorm:"not null"`
	AccountID           string              `gorm:"unique;not null"`
	LinkPassword        LinkPassword		`gorm:"embedded"`
	CompletedProgresses []CompletedProgress `gorm:"foreignKey:UserId" json:"completed_progresses"`
}
