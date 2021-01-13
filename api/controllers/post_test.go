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
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"gitlab.com/hydra/forum-api/api/controllers"
	"gitlab.com/hydra/forum-api/api/middlewares"
	"gitlab.com/hydra/forum-api/api/utils"
)

func Test_GetPost(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/posts").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.GetPosts)).
		Methods("GET")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success-Basic", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 5)).
			End()
	})

	t.Run("Success-Pagination", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Query("page", "2").
			Query("pagesize", "2").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 2)).
			Assert(jsonpath.Equal("S[0].id", float64(3))).
			Assert(jsonpath.Equal("S[1].id", float64(4))).
			End()
	})

	t.Run("Success-SelectedUserId", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Query("userid", "2").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 3)).
			Assert(jsonpath.Equal("S[0].id", float64(2))).
			Assert(jsonpath.Equal("S[1].id", float64(4))).
			Assert(jsonpath.Equal("S[2].id", float64(5))).
			End()
	})

	t.Run("Success-SelectedUsername", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Query("username", "naufal").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 1)).
			Assert(jsonpath.Equal("S[0].id", float64(1))).
			End()
	})

	t.Run("Success-Search", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Query("username", "hasan").
			Query("search", "is").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 1)).
			Assert(jsonpath.Equal("S[0].id", float64(5))).
			End()
	})

	t.Run("Fail-MalformedQuery", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Query("page", "1s").
			Query("pagesize", "1s").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-ConfictedQuery", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/posts").
			Header("Authorization", token).
			Query("username", "hasan").
			Query("userid", "1").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
}

func Test_CreatePost(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/posts").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.CreatePost)).
		Methods("POST")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"thread_id": 1,
				"content":   "i want to fly",
			}).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("Fail-EmptyContent", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"thread_id": 1,
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-FilledID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":        12,
				"thread_id": 1,
				"content":   "i want to fly",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-FilledInvalidAuthorID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"author_id": 2,
				"thread_id": 1,
				"content":   "i want to fly",
			}).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})
}

func Test_UpdatePost(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/posts").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.UpdatePost)).
		Methods("PATCH")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":      1,
				"content": "i want to fly",
			}).
			Expect(t).
			Status(http.StatusNoContent).
			End()
	})

	t.Run("Fail-UpdateOtherUserPost", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":      3,
				"content": "i want to fly",
			}).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Fail-MissingID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"content": "i want to fly",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-NonExisentID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":      12,
				"content": "i want to fly",
			}).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Fail-CantChangeAuthorID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/posts").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":        1,
				"content":   "i want to fly",
				"author_id": 2,
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func Test_DeletePost(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/posts/{id}").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.DeletePost)).
		Methods("DELETE")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Fail-DeleteOtherUserPost", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/posts/2").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Fail-MalformedIDPath", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/posts/1-").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-NonExisentID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/posts/12").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/posts/1").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Present("id")).
			End()
	})
}
