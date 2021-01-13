package controllers

import (
	// "fmt"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
	"gorm.io/gorm"
)

// GetPosts will fetch all posts of a specified criteria
func GetPosts(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	search := r.FormValue("search")
	userIDStr := r.FormValue("userid")
	pageNumStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pagesize")

	// fmt.Printf("%s, %s, %s, %s\n", username, userIDStr, pageNumStr, pageSizeStr)
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

	if !utils.IsPositiveInteger(int(userID64), pageNum, pageSize) {
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

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
	} else if userID != 0 || user.ID != 0 {
		search = "%!!.1?)(3}6]78[@;!%"
	} else {
		search = "%"
	}

	var posts []models.Post
	if userID != 0 || user.ID != 0 {
		err = db.Model(&models.Post{}).
			Where("creator_id = ? OR creator_id = ? OR content LIKE ?",
				userID, user.ID, search).
			Offset(offset).
			Limit(pageSize).
			Find(&posts).Error
	} else {
		err = db.Model(&models.Post{}).
			Where("content LIKE ?", search).
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
func CreatePost(w http.ResponseWriter, r *http.Request) {

}

// UpdatePost will update an existing Post
func UpdatePost(w http.ResponseWriter, r *http.Request) {

}

// DeletePost will delete an existing Post
func DeletePost(w http.ResponseWriter, r *http.Request) {

}
