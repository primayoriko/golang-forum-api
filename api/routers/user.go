package routers

import (
	// "github.com/hydra/forum-api/api/controllers"
	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/controllers"
)

func AddUserRoutes(router *mux.Router) error {
	var route *mux.Route

	route = router.HandleFunc("/users/register", controllers.Register).Methods("POST")

	if route == nil {

	}

	return nil
}
