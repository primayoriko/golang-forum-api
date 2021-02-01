package seeders

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/primayoriko/golang-forum-api/api/database"
	"github.com/primayoriko/golang-forum-api/migrations"
)

// SeedData is function to seed all tables data to database
func SeedData() error {
	// if err := migrations.DropTables(); err != nil {
	// 	return err
	// }

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

	return nil
}
