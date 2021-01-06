package main

import (
	"net/http"

	"gitlab.com/hydra/forum-api/api"
)

func main() {
	api.Run()
	// http.HandleFunc("/", hello)
	// fmt.Println("Server started")
	// log.Fatal(http.ListenAndServe(":8008", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"hello world!"}`))
}
