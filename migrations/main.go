package migrations

import (
	"fmt"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
)

// DropTables used to specified tables that exist in the database
func DropTables() error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	var tableNames = []string{
		"posts", "threads", "users"}

	for _, tableName := range tableNames {
		if err := db.Migrator().DropTable(tableName); err != nil {
			return err
		}
	}

	fmt.Println("Table(s) dropped")

	return nil
}

// MigrateModels used to migrate table of models that need to be created in the database
func MigrateModels() error {
	db, err := database.ConnectDB()
	if err != nil {
		// log.Fatalf("Error connecting to db, %v", err)
		return err
	}

	var dbModels []interface{} = []interface{}{
		&models.User{}, &models.Thread{}, &models.Post{}}

	for _, model := range dbModels {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	fmt.Println("Model(s) migrated")

	return nil
}
