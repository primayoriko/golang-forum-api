package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/auth"
	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
	"gorm.io/gorm"
)

// SignUp is function for create new User
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if err := db.Select("username", "email", "password").Create(&user).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusCreated, nil, nil)
}

// SignIn is method for get token for creds/auth
func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds auth.Credentials
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized,
				*(models.NewErrorResponse("wrong password/username")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	isTruePass := utils.CheckPasswordHash(creds.Password, user.Password)
	if !isTruePass {
		utils.JSONResponseWriter(&w, http.StatusUnauthorized,
			*(models.NewErrorResponse("wrong password/username")), nil)
		return
	}

	claims := &auth.Claims{
		ID:             user.ID,
		Username:       creds.Username,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		map[string]interface{}{"token": tokenString}, nil)
	return
}

// GetUsers is method for getting data for specified user(s) based on some criteria
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// db, err := database.ConnectDB()
	// if err != nil {
	// 	utils.JSONResponseWriter(&w, http.StatusInternalServerError,
	// 		*(models.NewErrorResponse(err.Error())), nil)
	// 	return
	// }
	mux.Vars(r)

	utils.JSONResponseWriter(&w, http.StatusOK, nil, nil)
}

// UpdateUser is for change it's own user data
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// db, err := database.ConnectDB()
	// if err != nil {
	// 	utils.JSONResponseWriter(&w, http.StatusInternalServerError,
	// 		*(models.NewErrorResponse(err.Error())), nil)
	// 	return
	// }
}

// DeleteUser is for delete it's own user account
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// var cred Credential
}
