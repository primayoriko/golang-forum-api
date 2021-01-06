package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"gitlab.com/hydra/forum-api/api"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		// if err := seeds.SeedData(); err != nil {
		// 	log.Fatalf("Seeding data encountered error, %v", err)
		// }
		api.Run()
	}

	// http.HandleFunc("/", hello)
	// fmt.Println("Server started")
	// log.Fatal(http.ListenAndServe(":8008", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"hello world!"}`))
}
