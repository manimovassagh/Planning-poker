// models/task.go

package models

import (
	"time"
)

type Task struct {
	ID             uint      `gorm:"primaryKey"`
	SessionID      uint      `gorm:"not null"` // Foreign key to Session
	TaskName       string    `gorm:"not null"`
	TaskDescription string   `gorm:"type:text"`
	CreatedAt      time.Time
	Status         string    `gorm:"default:'pending'"` // Can be "pending", "in_progress", "completed", etc.
}