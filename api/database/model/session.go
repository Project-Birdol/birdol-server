package model

import ("gorm.io/gorm")

type Session struct {
	gorm.Model
	SessionID	string	`gorm:"unique;not null"`
	DeviceID	string	`gorm:"unique;not null"`
	UserID		uint	`gorm:"unique;not null"`
}
