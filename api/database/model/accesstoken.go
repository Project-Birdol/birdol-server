package model

import (
	"time"

	"gorm.io/gorm"
)

type AccessToken struct {
	UserID			uint		`gorm:"primaryKey"`
	DeviceID		string		`gorm:"unique;not null"`
	Token			string		`gorm:"not null"`
	TokenUpdated	time.Time	`gorm:"not null"`
	CreatedAt		time.Time
	UpdatedAt		time.Time
	DeletedAt		gorm.DeletedAt	`gorm:"index"`
}

/* To disable soft delete
type AccessToken2 struct {
	Model
	UserID			uint		`gorm:"not null"`
	DeviceID		string		`gorm:"unique;not null"`
	Token			string		`gorm:"not null"`
	TokenUpdated	time.Time	`gorm:"not null"`
}
*/