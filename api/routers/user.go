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
		}, controllers.SignUp)).
		Methods("POST").Name("SignUp")

	router.HandleFunc("/signin",
		utils.ChainHandlerFuncs([]utils.Middleware{
			middlewares.Log,
		}, controllers.SignIn)).
		Methods("POST").Name("SignIn")

	// userRouter := router.PathPrefix("/users").Subrouter()
	// // userRouter.HandleFunc("/{username:[a-zA-Z0-9]+}",
	// userRouter.Path("/{id}").
	// 	HandlerFunc(utils.ChainHandlerFuncs(
	// 		[]utils.Middleware{
	// 			middlewares.CheckJWT,
	// 			middlewares.Log,
	// 		}, controllers.GetUser)).
	// 	Methods("GET").Name("GetUser")

	// userRouter.Queries().
	// 	HandlerFunc(utils.ChainHandlerFuncs(
	// 		[]utils.Middleware{
	// 			middlewares.CheckJWT,
	// 			middlewares.Log,
	// 		}, controllers.GetUsers)).
	// 	Methods("GET").Name("GetUsers")

	// userRouter.Queries().
	// 	HandlerFunc(utils.ChainHandlerFuncs(
	// 		[]utils.Middleware{
	// 			middlewares.CheckJWT,
	// 			middlewares.Log,
	// 		}, controllers.UpdateUser)).
	// 	Methods("PATCH").Name("UpdateUser")

	// userRouter.Path("/{id}").
	// 	HandlerFunc(utils.ChainHandlerFuncs(
	// 		[]utils.Middleware{
	// 			middlewares.CheckJWT,
	// 			middlewares.Log,
	// 		}, controllers.DeleteUser)).
	// 	Methods("DELETE").Name("DeleteUser")

	return nil
}
