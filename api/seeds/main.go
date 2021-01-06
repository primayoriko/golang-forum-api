package seeds

import (
	"fmt"
	"log"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/migrations"
)

func SeedData() error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	if db != nil {
		fmt.Println("db loaded")
		if err := migrations.MigrateModels(db); err != nil {
			log.Fatalf("Error model migration occured, %v", err)
		} else {
			fmt.Println("Model migrated")
		}
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
