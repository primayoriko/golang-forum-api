package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/routers"
)

func Run() {
	// if err := migrations.MigrateModels(); err != nil {
	// 	log.Fatalf("Error model migration occured, %v", err)
	// } else {
	// 	fmt.Println("Models migrated")
	// }

	r := mux.NewRouter()

	routers.AddUserRoutes(r)
	routers.AddThreadRoutes(r)
	routers.AddPostRoutes(r)
	// http.Handle("/", r)

	fmt.Printf("Start server at localhost:%s!\n", os.Getenv("API_PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("API_PORT")), r))
}
