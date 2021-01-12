package models

import (
	"time"
)

// ThreadResult is model for threads in the response
type ThreadResult struct {
	ID        uint64    `gorm:"serial" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Topic     string    `gorm:"size:255;not null" json:"topic"`
	CreatorID uint32    `gorm:"not null" json:"creator_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
