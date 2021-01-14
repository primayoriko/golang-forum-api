package models

import (
	"errors"
	"time"

	"gitlab.com/hydra/forum-api/api/utils"
)

// Post is model for posts table in the db
type Post struct {
	ID       uint64  `gorm:"serial" json:"id"`
	Author   *User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"author"`
	AuthorID uint32  `gorm:"not null" json:"author_id"`
	Thread   *Thread `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"thread"`
	ThreadID uint64  `gorm:"not null" json:"thread_id"`
	// ThreadType string
	Content   string    `gorm:"text;not null;" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// PostCreateRequest is request body schema for creating thread
type PostCreateRequest struct {
	ThreadID uint64 `json:"thread_id"`
	Content  string `json:"content"`
}

// PostUpdateRequest is request body schema for updating thread
type PostUpdateRequest struct {
	ID      uint64 `json:"id"`
	Content string `json:"content"`
}

// InjectToModel method to Injecting request structure to it's coresponsing model
func (pc *PostCreateRequest) InjectToModel(target *Post) error {
	// if pc.ThreadID == 0 || pc.Content == "" {
	if !utils.IsNonEmpty(pc.ThreadID, pc.Content) {
		return errors.New("thread_id/content cannot be left blank/empty")
	}

	target.ThreadID = pc.ThreadID
	target.Content = pc.Content

	return nil
}

// InjectToModel method to Injecting request structure to it's coresponsing model
func (pu *PostUpdateRequest) InjectToModel(target *Post) error {
	if pu.ID == 0 {
		return errors.New("id/content cannot be left blank/empty")
	}

	if pu.Content != "" {
		target.Content = pu.Content
	}

	target.ID = pu.ID

	return nil
}
