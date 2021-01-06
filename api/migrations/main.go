package migrations

import (
	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/gorm"
)

func MigrateModels(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Thread{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Post{}); err != nil {
		return err
	}

	return nil
}
