package seeders

import (
	"errors"

	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/gorm"
)

var threads []models.Thread = []models.Thread{
	{ID: 2, CreatorID: 1, Topic: "CP", Title: "How to Win CP"},
	{ID: 4, CreatorID: 2, Topic: "CP", Title: "How to Win CP 3"},
	{ID: 5, CreatorID: 2, Topic: "Ibadah", Title: "Tawakal"},
}

// SeedThreads is function for seed Post data
func SeedThreads(db *gorm.DB) error {
	if tx := db.Create(&threads); tx.Error == nil {
		return errors.New("seed threads data failed")
	}
	return nil
}
