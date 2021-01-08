package seeders

import (
	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/gorm"
)

var posts []models.Post = []models.Post{
	models.Post{AuthorID: 1, ThreadID: 2, Content: "Good! This is the start"},
	models.Post{AuthorID: 2, ThreadID: 4, Content: "Good! This is the start"},
	models.Post{AuthorID: 3, ThreadID: 4, Content: "Good! This is the start"},
	models.Post{AuthorID: 2, ThreadID: 2, Content: "Good! This is the start"},
	models.Post{AuthorID: 2, ThreadID: 5, Content: "Good! This is the start"},
}

// SeedPosts is function for seed Post data
func SeedPosts(db *gorm.DB) error {
	if tx := db.Create(&posts); tx.Error == nil {
		return tx.Error
	}
	return nil
}
