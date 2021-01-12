package routers

import (
	// "github.com/hydra/forum-api/api/controllers"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/controllers"
	"gitlab.com/hydra/forum-api/api/middlewares"
	"gitlab.com/hydra/forum-api/api/utils"
)

// AddUserRoutes is function to add subroute for auth and /users prefixes path
func AddUserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.Log,
		}, controllers.SignUp)).Methods("POST")
	router.HandleFunc("/signin",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.Log,
		}, controllers.SignIn)).Methods("POST")

	userRouter := router.PathPrefix("/users").Subrouter()
	// userRouter.HandleFunc("/{username:[a-zA-Z0-9]+}",
	userRouter.HandleFunc("/{username}",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.GetUsers)).Methods("GET")

	userRouter.HandleFunc("/",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.GetUsers)).Methods("GET")

	userRouter.HandleFunc("/",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.UpdateUser)).Methods("PATCH")

	userRouter.HandleFunc("/{username}",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.CheckJWT,
			middlewares.Log,
		}, controllers.DeleteUser)).Methods("DELETE")

	return nil
}
