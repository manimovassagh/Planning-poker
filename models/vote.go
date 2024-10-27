// models/vote.go

package models

import (
	"time"
)

type Vote struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"` // Foreign key to Task
	UserID    uint      `gorm:"not null"` // Foreign key to User
	VoteValue int       `gorm:"not null"`
	CreatedAt time.Time
}