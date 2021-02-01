package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/primayoriko/golang-forum-api/api"
	"github.com/primayoriko/golang-forum-api/migrations"
	"github.com/primayoriko/golang-forum-api/seeders"
)

func main() {
	err := godotenv.Load(".env")

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
