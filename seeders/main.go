package seeders

import (
	"fmt"

	"gitlab.com/hydra/forum-api/database"
	"gitlab.com/hydra/forum-api/migrations"
)

// SeedData is function to seed all tables data to database
func SeedData() error {
	if err := migrations.MigrateModels(); err != nil {
		return err
	}

	fmt.Println("Models migrated")

	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	fmt.Println("User data created")

	if err := SeedThreads(db); err != nil {
		return err
	}

	fmt.Println("Thread data created")

	if err := SeedPosts(db); err != nil {
		return err
	}

	fmt.Println("Post data created")

	return nil
}
