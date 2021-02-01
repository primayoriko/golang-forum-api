package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/primayoriko/golang-forum-api/api/database"
	"github.com/primayoriko/golang-forum-api/api/models"
	"github.com/primayoriko/golang-forum-api/api/utils"
)

// GetThreads will fetch all threads list of specific criteria
// @Title Get Threads.
// @Description Get all related thread of specific criteria.
// @Param username query string	optional	"User.Username"
// @Param userid query uint32	optional	"Thread.CreatorID -> User.ID"
// @Param topic query string	optional	"Thread.Topic"
// @Param title query string	optional	"Thread.Title"
// @Param page query int	optional	"pagination, current page"
// @Param pagesize query int	optional	"pagination, entry per page"
// @Success  200  array  models.Thread  "Thread JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /threads [get]
// @Tag Thread
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

	if title != "" {
		title = fmt.Sprintf("%%%s%%", title)
	} else {
		title = "%"
	}

	var threads []models.Thread
	if userID != 0 || user.ID != 0 || topic != "" {
		err = db.Model(&models.Thread{}).
			Where("(creator_id = ? OR creator_id = ? OR topic =  ?) AND title LIKE ?",
				userID, user.ID, topic, title).
			Order("created_at desc").
			Offset(offset).
			Limit(pageSize).
			Find(&threads).Error
	} else {
		err = db.Model(&models.Thread{}).
			Where("title LIKE ?", title).
			Order("created_at desc").
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
// @Title Get Thread.
// @Description Get a thread by it's ID.
// @Param  id  path  int  true  "Thread.ID"
// @Success  200  object  models.Thread  "Thread JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /threads/{id} [get]
// @Tag Thread
func GetThread(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	if idStr == "" || !utils.IsInteger(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid id format")), nil)
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	var thread models.Thread
	var posts []models.Post
	// db.Preload("posts").Where("id = ?", id).First(&thread)
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

	if err := db.Where("thread_id = ?", id).Order("updated_at desc").Limit(10).
		Find(&posts).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	thread.Posts = posts

	utils.JSONResponseWriter(&w, http.StatusOK,
		thread, nil)
	return
}

// CreateThread will make a new thread
// @Title Create Thread.
// @Description Create a new thread from JSON-formatted request body.
// @Param  thread  body  models.ThreadCreateRequest  true  "ThreadCreateRequest"
// @Success  201  object  models.ErrorResponse   "Created - No Body"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /threads [post]
// @Tag Thread
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
// @Title Update Thread.
// @Description Update an existing thread from JSON-formatted request body.
// @Param  thread  body  models.ThreadUpdateRequest  true  "ThreadUpdateRequest"
// @Success  204  object  models.ErrorResponse   "No Content - No Body"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /threads [patch]
// @Tag Thread
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

	if err := db.Model(&thread).Updates(thread).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent,
		nil, nil)
	return
}

// DeleteThread will delete an existing Thread
// @Title Delete Thread.
// @Description Delete an existing thread by it's ID.
// @Param  id  path  int  true  "Thread.ID"
// @Success  200  object  models.Thread  "Thread JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /threads/{id} [delete]
// @Tag Thread
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
