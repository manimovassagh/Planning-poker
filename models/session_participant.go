package models

import (
	"time"
)

type SessionParticipant struct {
	ID        uint      `gorm:"primaryKey"`
	SessionID uint      `gorm:"not null"` // Foreign key to Session
	UserID    uint      `gorm:"not null"` // Foreign key to User
	Role      string    `gorm:"type:varchar(20);default:'participant'"` // "admin" or "participant"
	JoinedAt  time.Time
}