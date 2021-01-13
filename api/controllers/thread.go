package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
)

// GetThreads will fetch all threads list of specific criteria
// @Title Get Thread List.
// @Description Get all related thread of specific criteria.
// @Param  groupID  path  int  true  "Id of a specific group."
// @Success  200  array  models.Thread  "Thread JSON"
// @Resource threads
// @Route /threads [get]
func GetThreads(w http.ResponseWriter, r *http.Request) {
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
					*(models.NewErrorResponse("can't find any threads")), nil)
				return
			}

			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				*(models.NewErrorResponse(err.Error())), nil)
			return
		}

		if userID != 0 && userID != user.ID {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find any threads")), nil)
			return
		}
	}

	var threads []models.Thread
	if userID != 0 || user.ID != 0 || topic != "" || title != "" {
		err = db.Model(&models.Thread{}).
			Where("creator_id = ? OR creator_id = ? OR topic =  ? OR title = ?",
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
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		threads, nil)
	return
}

// GetThread will fetch threads and its posts list of specific criteria
func GetThread(w http.ResponseWriter, r *http.Request) {
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

	var thread models.Thread
	// db.Preload("posts").
	// 	Where("id = ?", id).
	// 	First(&thread)

	err = db.Where("id = ?", id).
		First(&thread).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified thread")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		thread, nil)
	return
}

// CreateThread will make a new thread
func CreateThread(w http.ResponseWriter, r *http.Request) {
	var threadReq models.ThreadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&threadReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var thread models.Thread
	if err := threadReq.InjectToModel(&thread); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	thread.CreatorID = context.Get(r, "id").(uint32)

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	if err := db.Select("title", "topic", "creator_id").
		Create(&thread).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusCreated, nil, nil)
}

// UpdateThread will update an existing Thread
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)

	var threadReq models.ThreadUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&threadReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var thread, dbThread models.Thread
	if err := threadReq.InjectToModel(&thread); err != nil {
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

	if err := db.Where("id = ?", thread.ID).First(&dbThread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified thread")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if dbThread.CreatorID != userID {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do specified action as this user")), nil)
		return
	}

	if thread.Topic != "" {
		dbThread.Topic = thread.Topic
	}
	if thread.Title != "" {
		dbThread.Title = thread.Title
	}

	if err := db.Model(&thread).Updates(dbThread).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent,
		nil, nil)
	return
}

// DeleteThread will delete an existing Thread
func DeleteThread(w http.ResponseWriter, r *http.Request) {
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

	var thread models.Thread
	if err := db.Where("id = ?", id).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified thread")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)

		return
	}

	if thread.CreatorID != userID.(uint32) {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do the action as this user")), nil)
		return
	}

	if err := db.Delete(&thread, id).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)

		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		thread, nil)
	return
}
