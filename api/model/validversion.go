package model

type ValidClient struct {
	Model
	Platform string `gorm:"not null;unique"`
	SystemVersion uint `gorm:"not null"`
	MajorVersion uint `gorm:"not null"`
	MinorVersion uint `gorm:"not null"`
	Build string `gorm:"not null"`
}
