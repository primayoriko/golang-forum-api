package seeders

import (
	"fmt"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/migrations"
)

func SeedData() error {
	if err := migrations.MigrateModels(); err != nil {
		return err
	} else {
		fmt.Println("Models migrated")
	}

	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	} else {
		fmt.Println("User data created")
	}

	if err := SeedThreads(db); err != nil {
		return err
	} else {
		fmt.Println("Thread data created")
	}

	if err := SeedPosts(db); err != nil {
		return err
	} else {
		fmt.Println("Post data created")
	}

	return nil
}
