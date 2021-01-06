package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
	// "github.com/dgrijalva/jwt-go"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func register(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userData.Password, err = utils.HashPassword(userData.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func signIn(w http.ResponseWriter, r *http.Request) {
	// var cred Credential
}
