package migrations

import (
	"fmt"

	"github.com/primayoriko/golang-forum-api/api/database"
	"github.com/primayoriko/golang-forum-api/api/models"
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

// MigrateModels used to migrate table of models that need to be created in the database by dropping existing tables then re-create it
func MigrateModels() error {
	if err := DropTables(); err != nil {
		return err
	}

	db, err := database.ConnectDB()
	if err != nil {
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
