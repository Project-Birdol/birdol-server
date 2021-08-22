package model

import "time"

// To disable soft delete
type Model struct {
    ID        uint		`gorm:"primarykey"`  
    CreatedAt time.Time  
    UpdatedAt time.Time
}