package prisma

import (
	"time"
)

// Alias model for mapping filenames to aliases
type Alias struct {
	ID       uint   `gorm:"primaryKey"`
	FileName string `gorm:"uniqueIndex"`
	Alias    string `gorm:"uniqueIndex"`
	Port     int
	IsActive bool `gorm:"default:false"`
}

// SystemConfig model for storing system configuration
type SystemConfig struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"uniqueIndex"`
	Value       string
	Type        string    `gorm:"default:string"` // string, number, boolean, json
	Description string    `gorm:"default:''"`     // optional description
	HideValue   bool      `gorm:"default:false"`  // hide value in the UI
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
