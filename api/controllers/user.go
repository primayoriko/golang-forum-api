package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
)

// Register is function for create new User
func Register(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	password, ok := user["password"].(string)
	if !ok {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	user["password"], err = utils.HashPassword(password)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	fmt.Println(user)

	if tx := db.Model(&models.User{}).Create(&user); tx.Error != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusCreated, nil, nil)
}

// SignIn is method for get token for creds/auth
func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	if tx := db.Where("username = ?", creds.Username).First(&user); tx.Error != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	isTruePass := utils.CheckPasswordHash(creds.Password, user.Password)
	if !isTruePass {
		utils.JSONResponseWriter(&w, http.StatusUnauthorized,
			map[string]interface{}{
				"message": "wrong password/username"}, nil)
		return
	}

	claims := &models.Claims{
		Username:       creds.Username,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		map[string]interface{}{
			"token": tokenString}, nil)
	return
}

// ChangeUserData is for change it's own user data
func ChangeUserData(w http.ResponseWriter, r *http.Request) {
	// var cred Credential
}

// DeleteUserAccount is for delete it's own user account
func DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	// var cred Credential
}
