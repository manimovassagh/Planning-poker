// models/session.go

package models

import (
	"time"
)
type Session struct {
	ID          uint      `gorm:"primaryKey"`
	SessionName string    `gorm:"not null"`
	CreatedBy   uint      `gorm:"not null"` // Foreign key to User (admin ID)
	AdminID     uint      `gorm:"not null"`
	CreatedAt   time.Time
	IsActive    bool      `gorm:"default:true"` // Becomes false when closed
}