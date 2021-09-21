package model

import "time"

type AccessToken struct {
	Model
	UserID			uint		`gorm:"not null"`
	DeviceID		string		`gorm:"unique;not null"`
	Token			string		`gorm:"unique;not null"`
	RefreshToken	string		`gorm:"not null"`
	TokenUpdated	time.Time	`gorm:"not null"`
	PublicKey		string		`gorm:"not null"`
}
