package seeders

import (
	"github.com/primayoriko/golang-forum-api/api/models"
	"gorm.io/gorm"
)

var threads []models.Thread = []models.Thread{
	models.Thread{CreatorID: 1, Topic: "CP", Title: "How to Win CP"},
	models.Thread{CreatorID: 2, Topic: "CP", Title: "How to Win CP 3"},
	models.Thread{CreatorID: 2, Topic: "Ibadah", Title: "Tawakal"},
}

// SeedThreads is function for seed Post data
func SeedThreads(db *gorm.DB) error {
	if result := db.Create(&threads); result.Error == nil {
		return result.Error
	}
	return nil
}
