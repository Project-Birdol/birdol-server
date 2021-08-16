package model

import (
	"time"

	"gorm.io/gorm"
)

type AccessToken struct {
	UserID       uint      `gorm:"primaryKey"`
	Token        string    `gorm:"not null"`
	TokenUpdated time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
