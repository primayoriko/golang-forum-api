package controllers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"gitlab.com/hydra/forum-api/api/controllers"
)

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

	t.Run("Failure", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/signup").
			JSON(map[string]string{
				"username": "test1",
				"email":    "t@g.com",
				"password": "123",
			}).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

func Test_Login(t *testing.T) {
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

	t.Run("Failure-WrongPass", func(t *testing.T) {
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

	t.Run("Failure-NotExistUser", func(t *testing.T) {
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
