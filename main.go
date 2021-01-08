package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/hydra/forum-api/api"
	"gitlab.com/hydra/forum-api/migrations"
	"gitlab.com/hydra/forum-api/seeders"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		if len(os.Args) > 2 && os.Args[2] == "seed" {
			if err := seeders.SeedData(); err != nil {
				log.Fatalf("Error encountered when seed data, %v", err)
			}
		} else if len(os.Args) > 2 && os.Args[2] == "migrate" {
			if err := migrations.MigrateModels(); err != nil {
				log.Fatalf("Error encountered when migrate models, %v", err)
			}
		} else {
			api.Run()
		}
	}
}
