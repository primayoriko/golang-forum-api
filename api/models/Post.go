package models

import (
	"time"
)

// Post is model for posts table in the db
type Post struct {
	ID        uint64    `gorm:"serial" json:"id"`
	Author    *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	Thread    *Thread   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"thread"`
	ThreadID  uint64    `gorm:"not null" json:"thread_id"`
	Content   string    `gorm:"text;not null;" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
