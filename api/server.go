package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MigrateModels(db *gorm.DB) error {
	var err error
	err = db.AutoMigrate(&models.User{})
	err = db.AutoMigrate(&models.Thread{})
	err = db.AutoMigrate(&models.Post{})
	return err
}

func Run() {
	var err error

	db, err := ConnectDB()

	if err != nil {
		log.Fatalf("Error connecting to db, %v", err)
	} else {
		fmt.Println("Connected to db")
	}

	if db != nil {
		fmt.Println("db loaded")
		err = MigrateModels(db)
		if err != nil {
			log.Fatalf("Error migrating model, %v", err)
		} else {
			fmt.Println("Model migrated")
		}
	}

	r := mux.NewRouter()

	// http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("API_PORT")), r))
}
