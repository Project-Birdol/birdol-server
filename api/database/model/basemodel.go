package model

import "time"

// To disable soft delete
type Model struct {
    ID        uint		`gorm:"primarykey" json:"id"`  
    CreatedAt time.Time  `json:"-"` 
    UpdatedAt time.Time  `json:"-"`
}