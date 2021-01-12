package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	validate "github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
)

// GetThreadsList will fetch all threads list of specific criteria
func GetThreadsList(w http.ResponseWriter, r *http.Request) {
	// cUsername := context.Get(r, "username")
	// cUserID := context.Get(r, "id")
	username := r.FormValue("username")
	topic := r.FormValue("topic")
	title := r.FormValue("title")
	userIDStr := r.FormValue("userid")
	pageNumStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pagesize")

	// fmt.Printf("%s, %s, %s, %s\n", username, userIDStr, pageNumStr, pageSizeStr)

	if !validate.IsInt(userIDStr) || !validate.IsInt(pageNumStr) || !validate.IsInt(pageSizeStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			map[string]interface{}{"message": "bad query value"}, nil)
		return
	}

	userID64, _ := strconv.ParseUint(userIDStr, 10, 32)
	userID := uint32(userID64)
	pageNum, _ := strconv.Atoi(r.FormValue("page"))
	pageSize, _ := strconv.Atoi(r.FormValue("pagesize"))
	offset := pageNum * pageSize

	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": "failed to connect db"}, nil)
		return
	}

	var user models.User
	if username != "" {
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.JSONResponseWriter(&w, http.StatusNotFound,
					nil, nil)
				return
			}

			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				map[string]interface{}{"message": err}, nil)
			return
		}

		if userID != 0 && userID != user.ID {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				nil, nil)
			return
		}
	}

	var threads []models.Thread
	if userID != 0 || user.ID != 0 || topic != "" || title != "" {
		err = db.Model(&models.Thread{}).
			Where("ID = ? OR ID = ? OR topic =  ? OR title = ?",
				userID, user.ID, topic, title).
			Offset(offset).
			Limit(pageSize).
			Find(&threads).Error
	} else {
		err = db.Model(&models.Thread{}).
			Offset(offset).
			Limit(pageSize).
			Find(&threads).Error
	}

	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": err}, nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		interface{}(threads), nil)
	return
}

// GetThread will fetch threads and its posts list of specific criteria
func GetThread(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	if idStr == "" || !validate.IsInt(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			map[string]interface{}{"message": "invalid id format"}, nil)
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	// fmt.Printf("%s %d", idStr, id)

	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": "failed to connect db"}, nil)
		return
	}

	var thread models.Thread
	// db.Preload("posts").
	// 	Where("id = ?", id).
	// 	First(&thread)

	err = db.Where("id = ?", id).
		First(&thread).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				nil, nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": err}, nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		interface{}(thread), nil)
	return
}

// CreateThread will make a new thread
func CreateThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			map[string]interface{}{"message": "invalid body format"}, nil)
		return
	}

	thread.CreatorID = context.Get(r, "id").(uint32)

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": "failed to connect db"}, nil)
		return
	}

	if err := db.Select("title", "topic", "creator_id").Create(&thread).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": err}, nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusCreated, nil, nil)
}

// UpdateThread will update an existing Thread
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			map[string]interface{}{"message": "invalid body format"}, nil)
		return
	}

	thread.CreatorID = context.Get(r, "id").(uint32)

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": "failed to connect db"}, nil)
		return
	}

	if err := db.Where("id = ?", id).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				nil, nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": err}, nil)

		return
	}

}

// DeleteThread will delete an existing Thread
func DeleteThread(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id")
	idStr := mux.Vars(r)["id"]

	if idStr == "" || !validate.IsInt(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			map[string]interface{}{"message": "invalid id format"}, nil)
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": "failed to connect db"}, nil)
		return
	}

	var thread models.Thread
	if err := db.Where("id = ?", id).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				nil, nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": err}, nil)

		return
	}

	if thread.CreatorID != userID.(uint32) {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			nil, nil)
		return
	}

	if err := db.Delete(&thread, id).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": err}, nil)

		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		interface{}(thread), nil)
	return
}
