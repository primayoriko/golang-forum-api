package models

import (
	"time"
)

// Thread is model for threads table in the db
type Thread struct {
	ID        uint64    `gorm:"serial" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Topic     string    `gorm:"size:255;not null" json:"topic"`
	Creator   *User     `json:"creator"`
	CreatorID uint32    `gorm:"not null" json:"creator_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
