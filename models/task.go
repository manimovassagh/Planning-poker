package models

import (
	"time"
)

type Task struct {
	ID           uint      `gorm:"primaryKey"`
	SessionID    uint      `gorm:"not null"` // Foreign key to Session
	TaskName     string    `gorm:"not null"`
	TaskDescription string
	CreatedAt    time.Time
	Status       string    `gorm:"default:'pending'"`
	Votes        []Vote    `gorm:"foreignKey:TaskID"`
}