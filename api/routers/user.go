package routers

import (
	// "github.com/hydra/forum-api/api/controllers"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/controllers"
)

// AddUserRoutes is function to add subroute for /users prefixes path
func AddUserRoutes(router *mux.Router) error {
	var route *mux.Route

	router.HandleFunc("/signup", controllers.SignUp)
	router.HandleFunc("/signin", controllers.SignIn)

	// userRouters := router.PathPrefix("/users").Subrouter()
	// userRouters.use()
	// userRouters.HandleFunc("/signup", utils.PipelineHandlerFuncs([]utils.Middleware{}, controllers.Register)).Methods("POST")

	// middlewares.JwtCheck(controllers.Register)

	if route == nil {

	}

	return nil
}
