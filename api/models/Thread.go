package models

import (
	"time"
)

// Thread is model for threads table in the db
type Thread struct {
	ID        uint64 `gorm:"serial" json:"id"`
	Title     string `gorm:"size:255;not null" json:"title"`
	Topic     string `gorm:"size:255;not null" json:"topic"`
	Creator   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"creator"`
	CreatorID uint32 `gorm:"not null" json:"creator_id"`
	// Posts     []Post    `gorm:"polymorphic:Thread;polymorphicValue:threads;" json:"posts"`
	Posts     []Post    `gorm:"foreign_key:thread_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"posts"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// // ThreadCreateRequest is request body schema for creating thread
// type ThreadCreateRequest struct {
// 	Title string `gorm:"size:255;not null" json:"title"`
// 	Topic string `gorm:"size:255;not null" json:"topic"`
// }

// // ThreadUpdateRequest is request body schema for updating thread
// type ThreadUpdateRequest struct {
// 	ID    uint64 `gorm:"serial" json:"id"`
// 	Title string `gorm:"size:255;not null" json:"title"`
// 	Topic string `gorm:"size:255;not null" json:"topic"`
// }
