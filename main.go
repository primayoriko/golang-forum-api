package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/hydra/forum-api/api"
	"gitlab.com/hydra/forum-api/api/seeders"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		if len(os.Args) > 2 && os.Args[2] == "seed" {
			if err := seeders.SeedData(); err != nil {
				log.Fatalf("Error encountered when seed data, %v", err)
			} else {
				fmt.Println("Seed data successful!")
			}
		} else if len(os.Args) > 2 && os.Args[2] == "hello" {
			http.HandleFunc("/", hello)
			fmt.Println("Server started!")
			log.Fatal(http.ListenAndServe(":8008", nil))
		} else {
			api.Run()
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"hello world!"}`))
}
