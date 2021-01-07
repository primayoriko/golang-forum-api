package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
)

// Credential is data used for login
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claim is login result as an auth token
type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Register is function for create new User
func Register(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.JsonResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	password, ok := user["password"].(string)
	if !ok {
		utils.JsonResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	user["password"], err = utils.HashPassword(password)
	if err != nil {
		utils.JsonResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JsonResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	if tx := db.Model(&models.User{}).Create(&user); tx.Error != nil {
		utils.JsonResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.JsonResponseWriter(&w, http.StatusCreated, nil, nil)
}

// SignIn is method for get token for creds/auth
func SignIn(w http.ResponseWriter, r *http.Request) {
	var cred Credential

	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		utils.JsonResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JsonResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	// isTruePass := utils.CheckPasswordHash(cred.Password, )
}

// ChangeUserData is for change it's own user data
func ChangeUserData(w http.ResponseWriter, r *http.Request) {
	// var cred Credential
}

// DeleteUserAccount is for delete it's own user account
func DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	// var cred Credential
}
