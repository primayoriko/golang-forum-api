package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"gitlab.com/hydra/forum-api/api/auth"
	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
)

// SignUp is function for create new User
// @Title Sign Up.
// @Description Create a new user from JSON-formatted request body.
// @Param  post  body  models.RegistrationRequest  true  "RegistrationRequest"
// @Success  201  object  models.ErrorResponse   "Created - No Body"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /signup [post]
// @Tag User
func SignUp(w http.ResponseWriter, r *http.Request) {
	var regisReq models.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&regisReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var user, dbUser models.User
	if err := regisReq.InjectToModel(&user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var err error
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

	dbUser = models.User{}
	err = db.Where("username = ? OR email = ?", user.Username, user.Email).First(&dbUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	} else if err == nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("existing username/email")), nil)
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
// @Title Sign In.
// @Description Login with JSON-formatted request body.
// @Param  post  body  auth.Credentials  true  "auth.Credentials"
// @Success  200  object  auth.Claims   "auth.Claims"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /signin [post]
// @Tag User
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
// @Title Get Users.
// @Description Get all related user of specific criteria.
// @Param username query string	optional	"User.Username"
// @Param minid query uint32	optional	"min inclusive value of User.ID"
// @Param maxid query uint32	optional	"max incluseive value of User.ID"
// @Param page query int	optional	"pagination, current page"
// @Param pagesize query int	optional	"pagination, entry per page"
// @Success  200  array  models.User  "User JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /users [get]
// @Tag User
func GetUsers(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	pageNumStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pagesize")
	minIDStr := r.FormValue("minid")
	maxIDStr := r.FormValue("maxid")

	if !utils.IsInteger(minIDStr, maxIDStr, pageNumStr, pageSizeStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("bad query value")), nil)
		return
	}

	minID64, _ := strconv.ParseUint(minIDStr, 10, 32)
	maxID64, _ := strconv.ParseUint(maxIDStr, 10, 32)
	minID := uint32(minID64)
	maxID := uint32(maxID64)
	pageNum, _ := strconv.Atoi(r.FormValue("page"))
	pageNum--
	pageSize, _ := strconv.Atoi(r.FormValue("pagesize"))
	offset := pageNum * pageSize

	if !utils.IsNonNegative(int(minID),
		int(maxID), pageNum+1, pageSize) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("bad query value")), nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	if username != "" {
		username = fmt.Sprintf("%%%s%%", username)
	} else {
		username = "%"
	}

	if maxID == 0 {
		maxID = 2147483647*2 - 2
	}

	var users []models.User
	err = db.Model(&models.User{}).
		Where("username LIKE ? AND id BETWEEN ? AND ?",
			username, minID, maxID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var usersRes []models.UserResponse
	var userRes models.UserResponse
	for _, dbUser := range users {
		if err := userRes.InsertFromModel(dbUser); err != nil {
			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				*(models.NewErrorResponse(err.Error())), nil)
			return
		}

		usersRes = append(usersRes, userRes)
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		usersRes, nil)
	return
}

// GetUser is method for getting data for specified user(s) based on some criteria
// @Title Get User.
// @Description Get a user by it's ID.
// @Param  id  path  int  true  "User.ID"
// @Success  200  object  models.User  "User JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /threads/{id} [get]
// @Tag Thread
func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	if idStr == "" || !utils.IsInteger(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid id format")), nil)
		return
	}

	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint32(id64)

	// fmt.Printf("%s %d", idStr, id)

	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var userRes models.UserResponse
	if err := userRes.InsertFromModel(user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		userRes, nil)
	return
}

// UpdateUser is for change it's own user data
// @Title Update User.
// @Description Update an existing user from JSON-formatted Request Body.
// @Param  user  body  models.UserUpdateRequest  true  "UserUpdateRequest"
// @Success  204  object  models.ErrorResponse   "No Content - No Body"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /users [patch]
// @Tag User
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)

	var userReq models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var user, dbUser models.User
	if err := userReq.InjectToModel(&user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	if err := db.Where("id = ?", user.ID).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if dbUser.ID != userID {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do specified action as this user")), nil)
		return
	}

	dbUser = models.User{}
	err = db.Where("email = ? AND id <> ?", user.Email, user.ID).First(&dbUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	} else if err == nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("existing email")), nil)
		return
	}

	fmt.Println(user)
	fmt.Println(dbUser)
	fmt.Println(err)

	if user.Password != "" {
		if user.Password, err = utils.HashPassword(user.Password); err != nil {
			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				*(models.NewErrorResponse(err.Error())), nil)
			return
		}
	}

	if err := db.Model(&user).Updates(user).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent,
		nil, nil)
	return
}

// DeleteUser is for delete it's own user account
// @Title Delete User.
// @Description Delete an existing user by it's id.
// @Param  id  path  int  true  "User.ID"
// @Success  200  object  models.User  "User JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /users/{id} [delete]
// @Tag User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id")
	idStr := mux.Vars(r)["id"]

	if idStr == "" || !utils.IsInteger(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid id format")), nil)
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)

		return
	}

	if user.ID != userID.(uint32) {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do the action as this user")), nil)
		return
	}

	if err := db.Delete(&user, id).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)

		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		user, nil)
	return
}
