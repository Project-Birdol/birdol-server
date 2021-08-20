package model

import (
	"time"

	"gorm.io/gorm"
)

// To disable soft delete
type Model struct {  
    ID        uint		`gorm:"primarykey"`  
    CreatedAt time.Time  
    UpdatedAt time.Time  
} 

type Session struct {
	Model
	SessionID	string	`gorm:"unique;not null"`
	DeviceID	string	`gorm:"unique;not null"`
	UserID		uint	`gorm:"not null"`
	Disconnect	bool
}
