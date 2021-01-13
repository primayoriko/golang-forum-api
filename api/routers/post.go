package routers

import (
	// "github.com/hydra/forum-api/api/controllers"
	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/controllers"
	"gitlab.com/hydra/forum-api/api/middlewares"
	"gitlab.com/hydra/forum-api/api/utils"
)

// AddPostRoutes is function to add subroute for /posts prefixes path
func AddPostRoutes(router *mux.Router) error {
	postRouter := router.PathPrefix("/posts").Subrouter()
	// postRouter.Path("/{id}").
	// 	HandlerFunc(utils.ChainHandlerFuncs(
	// 		[]utils.Middleware{
	// 			middlewares.CheckJWT,
	// 			middlewares.Log,
	// 		}, controllers.GetPost)).
	// 	Methods("GET").Name("GetPost")

	postRouter.Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.GetPosts)).
		Methods("GET").Name("GetPosts")

	postRouter.Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.CreatePost)).
		Methods("POST").Name("CreatePost")

	postRouter.Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.UpdatePost)).
		Methods("PATCH").Name("UpdatePost")

	postRouter.Path("/{id}").
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
				middlewares.Log,
			}, controllers.DeletePost)).
		Methods("DELETE").Name("DeletePost")
	return nil
}
