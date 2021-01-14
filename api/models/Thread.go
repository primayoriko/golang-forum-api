package models

import (
	"errors"
	"time"

	"gitlab.com/hydra/forum-api/api/utils"
)

// Thread is model for threads table in the db
type Thread struct {
	ID        uint64 `gorm:"serial" json:"id"`
	Title     string `gorm:"size:255;not null" json:"title"`
	Topic     string `gorm:"size:255;not null" json:"topic"`
	Creator   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"creator"`
	CreatorID uint32 `gorm:"" json:"creator_id"`
	// Posts     []Post    `gorm:"polymorphic:Thread;polymorphicValue:threads;" json:"posts"`
	Posts     []Post    `gorm:"foreign_key:thread_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"posts"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ThreadCreateRequest is request body schema for creating thread
type ThreadCreateRequest struct {
	Title string `json:"title"`
	Topic string `json:"topic"`
}

// ThreadUpdateRequest is request body schema for updating thread
type ThreadUpdateRequest struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Topic string `json:"topic"`
}

// InjectToModel method to Injecting request structure to it's coresponsing model
func (tc *ThreadCreateRequest) InjectToModel(target *Thread) error {
	// if tc.Title == "" || tc.Topic == "" {
	if !utils.IsNonEmpty(tc.Title, tc.Topic) {
		return errors.New("topic/content cannot be left blank/empty")
	}

	target.Title = tc.Title
	target.Topic = tc.Topic

	return nil
}

// InjectToModel method to Injecting request structure to it's coresponsing model
func (tu *ThreadUpdateRequest) InjectToModel(target *Thread) error {
	if tu.ID == 0 {
		return errors.New("id cannot be left blank/empty")
	}

	if tu.Topic != "" {
		target.Topic = tu.Topic
	}

	if tu.Title != "" {
		target.Title = tu.Title
	}

	target.ID = tu.ID

	return nil
}
