package seeders

import (
	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/gorm"
)

var posts []models.Post = []models.Post{
	models.Post{AuthorID: 1, ThreadID: 2, Content: "Good! This is the start"},
	models.Post{AuthorID: 2, ThreadID: 3, Content: "Bad! too Bad"},
	models.Post{AuthorID: 3, ThreadID: 3, Content: "Good! This is the start"},
	models.Post{AuthorID: 2, ThreadID: 2, Content: "Mayday! Mayday!"},
	models.Post{AuthorID: 2, ThreadID: 1, Content: "This is beginning my friend..."},
}

// SeedPosts is function for seed Post data
func SeedPosts(db *gorm.DB) error {
	if result := db.Create(&posts); result.Error == nil {
		return result.Error
	}
	return nil
}
