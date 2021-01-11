package routers

import (
	// "github.com/hydra/forum-api/api/controllers"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/controllers"
	"gitlab.com/hydra/forum-api/api/middlewares"
	"gitlab.com/hydra/forum-api/api/utils"
)

// AddUserRoutes is function to add subroute for /users prefixes path
func AddUserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.Log,
		}, controllers.SignUp)).Methods("POST")
	router.HandleFunc("/signin",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.Log,
		}, controllers.SignIn)).Methods("POST")

	userRouters := router.PathPrefix("/users").Subrouter()
	// userRouters.HandleFunc("/{username:[a-zA-Z0-9]+}",
	userRouters.HandleFunc("/{username}",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.GetUsers)).Methods("GET")

	userRouters.HandleFunc("/",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.GetUsers)).Methods("GET")

	userRouters.HandleFunc("/",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.UpdateUser)).Methods("PATCH")

	userRouters.HandleFunc("/",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.DeleteUser)).Methods("DELETE")

	return nil
}
