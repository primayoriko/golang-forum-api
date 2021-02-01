package routers

import (
	// "github.com/hydra/forum-api/api/controllers"
	"github.com/gorilla/mux"
	"github.com/primayoriko/golang-forum-api/api/controllers"
	"github.com/primayoriko/golang-forum-api/api/middlewares"
	"github.com/primayoriko/golang-forum-api/api/utils"
)

// AddThreadRoutes is function to add subroute for /threads prefixes path
func AddThreadRoutes(router *mux.Router) error {
	threadRouter := router.PathPrefix("/threads").Subrouter()
	threadRouter.Path("/{id}").
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.GetThread)).
		Methods("GET").Name("GetThread")

	threadRouter.Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.GetThreads)).
		Methods("GET").Name("GetThreads")

	threadRouter.Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.CreateThread)).
		Methods("POST").Name("CreateThread")

	threadRouter.Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.UpdateThread)).
		Methods("PATCH").Name("UpdateThread")

	threadRouter.Path("/{id}").
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.DeleteThread)).
		Methods("DELETE").Name("DeleteThread")

	return nil
}
