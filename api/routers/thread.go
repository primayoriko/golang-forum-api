package routers

import (
	// "github.com/hydra/forum-api/api/controllers"
	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/controllers"
	"gitlab.com/hydra/forum-api/api/middlewares"
	"gitlab.com/hydra/forum-api/api/utils"
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

	threadRouter.Path("/").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.GetThreadsList)).
		Methods("GET").Name("GetThreadsList")

	threadRouter.Path("/").
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
