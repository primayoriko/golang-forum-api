package controllers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"

	"gitlab.com/hydra/forum-api/api/controllers"
	"gitlab.com/hydra/forum-api/api/middlewares"
	"gitlab.com/hydra/forum-api/api/utils"
)

// func Test_GetThreadsList(t *testing.T) {
// 	if err := godotenv.Load("../../.env"); err != nil {
// 		log.Fatal("Error loading .env file", err)
// 	}

// 	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

// 	r := mux.NewRouter()
// 	r.Path("/threads").
// 		Queries().
// 		HandlerFunc(utils.ChainHandlerFuncs(
// 			[]utils.Middleware{
// 				middlewares.CheckJWT,
// 			}, controllers.GetThreadsList)).
// 		Methods("GET")

// 	ts := httptest.NewServer(r)

// 	defer ts.Close()

// 	t.Run("Success-Basic", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusOK).
// 			Assert(jsonpath.Len("$", 3)).
// 			End()
// 	})

// 	t.Run("Success-Pagination", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads").
// 			Header("Authorization", token).
// 			Query("page", "1").
// 			Query("pagesize", "1").
// 			Expect(t).
// 			Status(http.StatusOK).
// 			Assert(jsonpath.Len("$", 1)).
// 			Assert(jsonpath.Equal("S[0].id", float64(2))).
// 			End()
// 	})

// 	t.Run("Success-SelectedUserId", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads").
// 			Header("Authorization", token).
// 			Query("userid", "2").
// 			Expect(t).
// 			Status(http.StatusOK).
// 			Assert(jsonpath.Len("$", 2)).
// 			Assert(jsonpath.Equal("S[0].id", float64(2))).
// 			End()
// 	})

// 	t.Run("Success-SelectedUsername", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads").
// 			Header("Authorization", token).
// 			Query("username", "hasan").
// 			Expect(t).
// 			Status(http.StatusOK).
// 			Assert(jsonpath.Len("$", 2)).
// 			Assert(jsonpath.Equal("S[0].id", float64(2))).
// 			End()
// 	})

// 	t.Run("Fail-MalformedQuery", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads").
// 			Header("Authorization", token).
// 			Query("userid", "1s").
// 			Expect(t).
// 			Status(http.StatusBadRequest).
// 			End()
// 	})

// 	t.Run("Fail-ConfictedQuery", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads").
// 			Header("Authorization", token).
// 			Query("username", "hasan").
// 			Query("userid", "1").
// 			Expect(t).
// 			Status(http.StatusNotFound).
// 			End()
// 	})
// }

// func Test_GetThread(t *testing.T) {
// 	if err := godotenv.Load("../../.env"); err != nil {
// 		log.Fatal("Error loading .env file", err)
// 	}

// 	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

// 	r := mux.NewRouter()
// 	r.Path("/threads/{id}").
// 		Queries().
// 		HandlerFunc(utils.ChainHandlerFuncs(
// 			[]utils.Middleware{
// 				middlewares.CheckJWT,
// 			}, controllers.GetThread)).
// 		Methods("GET")

// 	ts := httptest.NewServer(r)

// 	defer ts.Close()

// 	t.Run("Success", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads/1").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusOK).
// 			Assert(jsonpath.Equal("$.id", float64(1))).
// 			End()
// 	})

// 	t.Run("Fail-NonExisentID", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads/12").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusNotFound).
// 			End()
// 	})

// 	t.Run("Fail-MalformedIDPath", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Get("/threads/1s").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusBadRequest).
// 			End()
// 	})
// }

// func Test_CreateThread(t *testing.T) {
// 	if err := godotenv.Load("../../.env"); err != nil {
// 		log.Fatal("Error loading .env file", err)
// 	}

// 	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

// 	r := mux.NewRouter()
// 	r.Path("/threads").
// 		Queries().
// 		HandlerFunc(utils.ChainHandlerFuncs(
// 			[]utils.Middleware{
// 				middlewares.CheckJWT,
// 			}, controllers.CreateThread)).
// 		Methods("POST")

// 	ts := httptest.NewServer(r)

// 	defer ts.Close()

// 	t.Run("Success", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Post("/threads").
// 			Header("Authorization", token).
// 			JSON(map[string]string{
// 				"topic": "lifestyle",
// 				"title": "fitness",
// 			}).
// 			Expect(t).
// 			Status(http.StatusCreated).
// 			End()
// 	})

// 	t.Run("Fail-FilledID", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Post("/threads").
// 			Header("Authorization", token).
// 			JSON(map[string]interface{}{
// 				"id":    12,
// 				"topic": "lifestyle",
// 				"title": "fitness",
// 			}).
// 			Expect(t).
// 			Status(http.StatusBadRequest).
// 			End()
// 	})

// 	t.Run("Fail-FilledInvalidCreatorID", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Post("/threads").
// 			Header("Authorization", token).
// 			JSON(map[string]interface{}{
// 				"creator_id": 2,
// 				"topic":      "lifestyle",
// 				"title":      "fitness",
// 			}).
// 			Expect(t).
// 			Status(http.StatusForbidden).
// 			End()
// 	})
// }

func Test_UpdateThread(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/threads").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.UpdateThread)).
		Methods("PATCH")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/threads").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    1,
				"topic": "lifestyler",
				"title": "fitnesses",
			}).
			Expect(t).
			Status(http.StatusNoContent).
			End()
	})

	t.Run("Fail-UpdateOtherUserThread", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/threads").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    3,
				"topic": "lifestyle",
				"title": "fitness",
			}).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Fail-MissingID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/threads").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"topic": "lifestyle",
				"title": "fitness",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-NonExisentID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/threads").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    12,
				"topic": "lifestyle",
				"title": "fitness",
			}).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Fail-CantChangeCreatorID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/threads").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":         1,
				"topic":      "lifestyle",
				"title":      "fitness",
				"creator_id": 2,
			}).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})
}

// func Test_DeleteThread(t *testing.T) {
// 	if err := godotenv.Load("../../.env"); err != nil {
// 		log.Fatal("Error loading .env file", err)
// 	}

// 	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

// 	r := mux.NewRouter()
// 	r.Path("/threads/{id}").
// 		Queries().
// 		HandlerFunc(utils.ChainHandlerFuncs(
// 			[]utils.Middleware{
// 				middlewares.CheckJWT,
// 			}, controllers.DeleteThread)).
// 		Methods("DELETE")

// 	ts := httptest.NewServer(r)

// 	defer ts.Close()

// 	t.Run("Fail-DeleteOtherUserThread", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Delete("/threads/2").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusForbidden).
// 			End()
// 	})

// 	t.Run("Fail-MalformedIDPath", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Delete("/threads/1-").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusBadRequest).
// 			End()
// 	})

// 	t.Run("Fail-NonExisentID", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Delete("/threads/12").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusNotFound).
// 			End()
// 	})

// 	t.Run("Success", func(t *testing.T) {
// 		apitest.New().
// 			Handler(r).
// 			Delete("/threads/1").
// 			Header("Authorization", token).
// 			Expect(t).
// 			Status(http.StatusOK).
// 			Assert(jsonpath.Present("id")).
// 			End()
// 	})
// }
