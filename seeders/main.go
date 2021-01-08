package seeders

import (
	"fmt"

	"gorm.io/gorm"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/migrations"
)

// SeedData is function to seed all tables data to database
func SeedData() error {
	if err := migrations.DropTables(); err != nil {
		return err
	}

	if err := migrations.MigrateModels(); err != nil {
		return err
	}

	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	seedersFunc := []func(*gorm.DB) error{
		SeedUsers, SeedThreads, SeedPosts}

	for _, function := range seedersFunc {
		if err := function(db); err != nil {
			return err
		}
	}

	fmt.Println("Data successfully seeded")

	// if err := SeedUsers(db); err != nil {
	// 	return err
	// }

	// fmt.Println("User data created")

	// if err := SeedThreads(db); err != nil {
	// 	return err
	// }

	// fmt.Println("Thread data created")

	// if err := SeedPosts(db); err != nil {
	// 	return err
	// }

	// fmt.Println("Post data created")

	return nil
}
