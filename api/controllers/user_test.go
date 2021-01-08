package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"

	"gitlab.com/hydra/forum-api/api/controllers"
)

func TestRegister(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.Register)

	ts := httptest.NewServer(r)

	defer ts.Close()

	apitest.New().
		Handler(r).
		Post("/register").
		JSON("").
		Expect(t).
		Status(http.StatusOK).
		End()

	// apitest.APITest
}

func TestLogin(t *testing.T) {
	r := mux.NewRouter()
}
