package migrations

import (
	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
)

func MigrateModels() error {
	db, err := database.ConnectDB()

	if err != nil {
		// log.Fatalf("Error connecting to db, %v", err)
		return err
	}

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
