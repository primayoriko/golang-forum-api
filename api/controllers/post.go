package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/primayoriko/golang-forum-api/api/database"
	"github.com/primayoriko/golang-forum-api/api/models"
	"github.com/primayoriko/golang-forum-api/api/utils"
	"gorm.io/gorm"
)

// GetPosts will fetch all posts of a specified criteria
// @Title Get Posts.
// @Description Get all related post of specific criteria.
// @Param username query string	optional	"User.Username"
// @Param userid query uint32	optional	"Post.AuthorID -> User.ID"
// @Param threadid query unit64	optional	"Post.ThreadID -> Thread.ID"
// @Param content query string	optional	"Post.Content"
// @Param page query int	optional	"pagination, current page"
// @Param pagesize query int	optional	"pagination, entry per page"
// @Success  200  array  models.Post  "Post JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /posts [get]
// @Tag Post
func GetPosts(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	content := r.FormValue("content")
	userIDStr := r.FormValue("userid")
	threadIDStr := r.FormValue("threadid")
	pageNumStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pagesize")

	// fmt.Printf("%s, %s, %s, %s\n", username, userIDStr, pageNumStr, pageSizeStr)
	if !utils.IsInteger(userIDStr, pageNumStr, pageSizeStr, threadIDStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("bad query value")), nil)
		return
	}

	threadID, _ := strconv.ParseUint(threadIDStr, 10, 64)
	userID64, _ := strconv.ParseUint(userIDStr, 10, 32)
	userID := uint32(userID64)
	pageNum, _ := strconv.Atoi(r.FormValue("page"))
	pageNum--
	pageSize, _ := strconv.Atoi(r.FormValue("pagesize"))
	offset := pageNum * pageSize

	if !utils.IsNonNegative(pageNum+1, pageSize, int(threadID)) {
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
					*(models.NewErrorResponse("can't find any posts")), nil)
				return
			}

			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				*(models.NewErrorResponse(err.Error())), nil)
			return
		}

		if userID != 0 && userID != user.ID {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find any posts")), nil)
			return
		}
	}

	if content != "" {
		content = fmt.Sprintf("%%%s%%", content)
	} else {
		content = "%"
	}

	var posts []models.Post
	if userID != 0 || user.ID != 0 || threadID != 0 {
		err = db.Model(&models.Post{}).
			Where("(author_id = ? OR author_id = ? OR thread_id = ?) AND content LIKE ?",
				userID, user.ID, threadID, content).
			Order("created_at desc").
			Offset(offset).
			Limit(pageSize).
			Find(&posts).Error
	} else {
		err = db.Model(&models.Post{}).
			Where("content LIKE ?", content).
			Order("created_at desc").
			Offset(offset).
			Limit(pageSize).
			Find(&posts).Error
	}

	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		posts, nil)
	return
}

// CreatePost will make a new post on a specific post
// @Title Create Post.
// @Description Create a new post from JSON-formatted request body.
// @Param  post  body  models.PostCreateRequest  true  "PostCreateRequest"
// @Success  201  object  models.ErrorResponse   "Created - No Body"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /posts [post]
// @Tag Post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var postReq models.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var post models.Post
	if err := postReq.InjectToModel(&post); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if post.AuthorID != context.Get(r, "id").(uint32) && post.AuthorID != 0 {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do the action as this user")), nil)
		return
	}

	post.AuthorID = context.Get(r, "id").(uint32)

	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse("failed to connect db")), nil)
		return
	}

	if err := db.Select("thread_id", "author_id", "content").
		Create(&post).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusCreated, nil, nil)
}

// UpdatePost will update an existing Post
// @Title Update Post.
// @Description Update an existing post from JSON-formatted request body.
// @Param  post  body  models.PostUpdateRequest  true  "PostUpdateRequest"
// @Success  204  object  models.ErrorResponse   "No Content - No Body"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /posts [patch]
// @Tag Post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)

	var postReq models.PostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			*(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var post, dbPost models.Post
	if err := postReq.InjectToModel(&post); err != nil {
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

	if err := db.Where("id = ?", post.ID).First(&dbPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified thread")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if dbPost.AuthorID != userID {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do specified action as this user")), nil)
		return
	}

	if err := db.Model(&post).Updates(post).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent,
		nil, nil)
	return
}

// DeletePost will delete an existing Post
// @Title Delete Post.
// @Description Delete an existing post by it's ID.
// @Param  id  path  int  true  "Post.ID"
// @Success  200  object  models.Post  "Post JSON"
// @Failure  400  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  401  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  403  object  models.ErrorResponse  "ErrorResponse JSON"
// @Failure  500  object  models.ErrorResponse  "ErrorResponse JSON"
// @Route /posts/{id} [delete]
// @Tag Post
func DeletePost(w http.ResponseWriter, r *http.Request) {
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

	var post models.Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find specified post")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)

		return
	}

	if post.AuthorID != userID.(uint32) {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do the action as this user")), nil)
		return
	}

	if err := db.Delete(&post, id).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)

		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		post, nil)
	return
}
