package models

import (
	"time"
)

type Session struct {
	ID         uint      `gorm:"primaryKey"`
	SessionName string   `gorm:"not null"`
	CreatedBy  uint      `gorm:"not null"` // Foreign key to User
	AdminID    uint      // Optional field to store the creator’s User ID
	CreatedAt  time.Time
	IsActive   bool      `gorm:"default:true"`
	Tasks      []Task    `gorm:"foreignKey:SessionID"`
	Participants []SessionParticipant `gorm:"foreignKey:SessionID"`
}