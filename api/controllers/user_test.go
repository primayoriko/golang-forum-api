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

	"github.com/primayoriko/golang-forum-api/api/controllers"
	"github.com/primayoriko/golang-forum-api/api/middlewares"
	"github.com/primayoriko/golang-forum-api/api/utils"
)

func Test_GetUsers(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/users").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.GetUsers)).
		Methods("GET")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success-Basic", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 4)).
			End()
	})

	t.Run("Success-Pagination", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users").
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

	t.Run("Success-SearchUsername", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users").
			Header("Authorization", token).
			Query("username", "n").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 3)).
			Assert(jsonpath.Equal("S[2].id", float64(4))).
			Assert(jsonpath.Equal("S[1].id", float64(2))).
			End()
	})

	t.Run("Success-SetMinMaxID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users").
			Header("Authorization", token).
			Query("minid", "2").
			Query("maxid", "3").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 2)).
			Assert(jsonpath.Equal("S[0].id", float64(2))).
			Assert(jsonpath.Equal("S[1].id", float64(3))).
			End()
	})

	t.Run("Fail-MalformedQuery", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users").
			Header("Authorization", token).
			Query("minid", "1s").
			Query("page", "1s").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func Test_GetUser(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/users/{id}").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.GetUser)).
		Methods("GET")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users/2").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Equal("$.id", float64(2))).
			Assert(jsonpath.Equal("$.username", "hasan")).
			End()
	})

	t.Run("Fail-NonExisentID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users/12").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Fail-MalformedIDPath", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Get("/users/1s").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func Test_SignUp(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/signup", controllers.SignUp).Methods("POST")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signup").
			JSON(map[string]string{
				"username": "test1",
				"email":    "test@g.com",
				"password": "123",
			}).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("Fail-ExistingUsernameOrEmail", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signup").
			JSON(map[string]string{
				"username": "test1",
				"email":    "test@g.com",
				"password": "123",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-EmptyRequiredField", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signup").
			JSON(map[string]string{
				"username": "test1",
				"password": "123",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-BadFormatField", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signup").
			JSON(map[string]string{
				"username": "test1",
				"email":    "a",
				"password": "123",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}

func Test_Signin(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/signin", controllers.SignIn).Methods("POST")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signin").
			JSON(map[string]string{
				"username": "test1",
				"password": "123",
			}).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Chain().
				Present("token").
				End()).
			End()
	})

	t.Run("Fail-WrongPass", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signin").
			JSON(map[string]string{
				"username": "test1",
				"password": "1234",
			}).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Fail-NonExistentUser", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signin").
			JSON(map[string]string{
				"username": "test2",
				"password": "1234",
			}).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
}

func Test_UpdateUser(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/users").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.UpdateUser)).
		Methods("PATCH")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/users").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    1,
				"email": "prim@g.com",
			}).
			Expect(t).
			Status(http.StatusNoContent).
			End()
	})

	t.Run("Fail-TakenEmail", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/users").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    1,
				"email": "t@g.com",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-UpdateOtherUserUser", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/users").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    2,
				"email": "masterhasan@g.com",
			}).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Fail-MissingID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/users").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"email": "masterhasan@g.com",
			}).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-NonExisentID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Patch("/users").
			Header("Authorization", token).
			JSON(map[string]interface{}{
				"id":    12,
				"email": "masterhasan@g.com",
			}).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
}

func Test_DeleteUser(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	token := "Bearer " + os.Getenv("TEST_USER1_TOKEN")

	r := mux.NewRouter()
	r.Path("/users/{id}").
		Queries().
		HandlerFunc(utils.ChainHandlerFuncs(
			[]utils.Middleware{
				middlewares.CheckJWT,
			}, controllers.DeleteUser)).
		Methods("DELETE")

	ts := httptest.NewServer(r)

	defer ts.Close()

	t.Run("Fail-DeleteOtherUserUser", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/users/2").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Fail-MalformedIDPath", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/users/1-").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Fail-NonExisentID", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/users/12").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Delete("/users/1").
			Header("Authorization", token).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Present("id")).
			End()
	})
}
