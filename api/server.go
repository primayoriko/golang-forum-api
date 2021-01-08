package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/routers"
)

// Run would start server for the api
func Run() {
	r := mux.NewRouter()

	routers.AddUserRoutes(r)
	routers.AddThreadRoutes(r)
	routers.AddPostRoutes(r)

	writeTimeout, _ := strconv.Atoi(os.Getenv("WRITE_TIMEOUT_SEC"))
	readTimeout, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT_SEC"))

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("localhost:%s", os.Getenv("API_PORT")),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	fmt.Printf("Start server at localhost:%s!\n", os.Getenv("API_PORT"))
	log.Fatalln(server.ListenAndServe())
}

/*
https://stackoverflow.com/questions/64768950/how-to-use-specific-middleware-for-specific-routes-in-a-get-subrouter-in-gorilla
*/
