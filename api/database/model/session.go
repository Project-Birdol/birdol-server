package model

type Session struct {
	Model
	SessionID	string	`gorm:"unique;not null"`
	AccessToken	string	`gorm:"unique;not null"`
	UserID		uint	`gorm:"not null"`
	Expired		bool
}
