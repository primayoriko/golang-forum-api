package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
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
	username := r.FormValue("username")
	topic := r.FormValue("topic")
	title := r.FormValue("title")
	userIDStr := r.FormValue("userid")
	pageNumStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pagesize")

	if !utils.IsInteger(userIDStr, pageNumStr, pageSizeStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("bad query value")), nil)
		return
	}

	userID64, _ := strconv.ParseUint(userIDStr, 10, 32)
	userID := uint32(userID64)
	pageNum, _ := strconv.Atoi(r.FormValue("page"))
	pageNum--
	pageSize, _ := strconv.Atoi(r.FormValue("pagesize"))
	offset := pageNum * pageSize

	if !utils.IsNonNegative(pageNum+1, pageSize) {
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

	var user models.User
	if username != "" {
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.JSONResponseWriter(&w, http.StatusNotFound,
					*(models.NewErrorResponse("can't find any users")), nil)
				return
			}

			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				*(models.NewErrorResponse(err.Error())), nil)
			return
		}

		if userID != 0 && userID != user.ID {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find any users")), nil)
			return
		}
	}

	var users []models.User
	if userID != 0 || user.ID != 0 || topic != "" || title != "" {
		err = db.Model(&models.User{}).
			Where("creator_id = ? OR creator_id = ? OR topic =  ? OR title = ?",
				userID, user.ID, topic, title).
			Offset(offset).
			Limit(pageSize).
			Find(&users).Error
	} else {
		err = db.Model(&models.User{}).
			Offset(offset).
			Limit(pageSize).
			Find(&users).Error
	}

	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		users, nil)
	return
}

// GetUser is method for getting data for specified user(s) based on some criteria
func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	if idStr == "" || !utils.IsInteger(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid id format")), nil)
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	// fmt.Printf("%s %d", idStr, id)

	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	var user models.User
	// db.Preload("posts").
	// 	Where("id = ?", id).
	// 	First(&user)

	err = db.Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		user, nil)
	return
}

// UpdateUser is for change it's own user data
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)
	var user, dbUser models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	if user.ID == 0 {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("need user id")), nil)
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

	if dbUser.CreatorID != userID {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do specified action as this user")), nil)
		return
	}

	if user.CreatorID != userID && user.CreatorID != 0 {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("can't change creator id")), nil)
		return
	}

	// user.CreatorID = userID
	if user.Topic != "" {
		dbUser.Topic = user.Topic
	}
	if user.Title != "" {
		dbUser.Title = user.Title
	}

	if err := db.Save(&dbUser).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if err := db.Model(&dbUser).Updates(dbUser).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent,
		nil, nil)
	return
}

// DeleteUser is for delete it's own user account
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

	if user.CreatorID != userID.(uint32) {
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
