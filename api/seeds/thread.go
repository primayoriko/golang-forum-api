package seeds

import (
	"errors"

	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/gorm"
)

var Threads []models.Thread = []models.Thread{
	{ID: 2, CreatorID: 1, Topic: "CP", Title: "How to Win CP"},
	{ID: 4, CreatorID: 2, Topic: "CP", Title: "How to Win CP 3"},
	{ID: 5, CreatorID: 2, Topic: "Ibadah", Title: "Tawakal"},
}

func SeedThreads(db *gorm.DB) error {
	if tx := db.Create(&Threads); tx == nil {
		return errors.New("create failed")
	}
	return nil
}
