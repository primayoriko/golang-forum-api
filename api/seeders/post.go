package seeders

import (
	"errors"

	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/gorm"
)

var posts []models.Post = []models.Post{
	{ID: 1, AuthorID: 1, ThreadID: 2, Content: "Good! This is the start"},
	{ID: 3, AuthorID: 2, ThreadID: 4, Content: "Good! This is the start"},
	{ID: 2, AuthorID: 3, ThreadID: 4, Content: "Good! This is the start"},
	{ID: 10, AuthorID: 2, ThreadID: 2, Content: "Good! This is the start"},
	{ID: 5, AuthorID: 2, ThreadID: 5, Content: "Good! This is the start"},
}

// SeedPosts is function for seed Post data
func SeedPosts(db *gorm.DB) error {
	if tx := db.Create(&posts); tx.Error == nil {
		return errors.New("seed posts data failed")
	}
	return nil
}
